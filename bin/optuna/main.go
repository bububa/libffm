package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path"
	"runtime"
	"sync"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/rdb"
	"github.com/c-bata/goptuna/tpe"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	dataPath     string
	testDataPath string
	libffmPath   string
	workPath     string
)

func init() {
	flag.StringVar(&dataPath, "tr", "", "data path")
	flag.StringVar(&testDataPath, "p", "", "validation data path")
	flag.StringVar(&libffmPath, "libffm", "ffm-train", "ffm-train file path")
	flag.StringVar(&workPath, "w", "", "work path")
	flag.Parse()
}

func main() {
	// setup storage
	dbPath := path.Join(workPath, "db.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("failed to open db:", err)
	}
	defer db.Close()
	db.DB().SetMaxOpenConns(1)
	rdb.RunAutoMigrate(db)
	storage := rdb.NewStorage(db)
	studyName := "goptuna-libffm"
	// create a study
	study, err := goptuna.CreateStudy(
		studyName,
		goptuna.StudyOptionStorage(storage),
		goptuna.StudyOptionSampler(tpe.NewSampler()),
	)
	if err != nil {
		log.Fatal("failed to create study:", err)
	}

	// create a context with cancel function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	study.WithContext(ctx)
	var wg sync.WaitGroup
	/*
		// set signal handler
		sigch := make(chan os.Signal, 1)
		defer close(sigch)
		signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		wg.Add(1)
		go func() {
			defer wg.Done()
			sig, ok := <-sigch
			if !ok {
				return
			}
			cancel()
			log.Print("catch a kill signal:", sig.String())
		}()
	*/
	// run optimize with context
	concurrency := runtime.NumCPU() - 1
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := study.Optimize(objective, 1000/concurrency)
			if err != nil {
				log.Print("optimize catch error:", err)
			}
		}()
	}
	wg.Wait()

	// print best hyper-parameters and the result
	v, _ := study.GetBestValue()
	params, _ := study.GetBestParams()
	log.Printf("Best evaluation=%f (lambda=%v, eta=%v, latent=%v)",
		v, params["lambda"], params["eta"], params["latent"])
}

func objective(trial goptuna.Trial) (float64, error) {
	lmd, err := trial.SuggestLogUniform("lambda", 1e-6, 1)
	if err != nil {
		return -1, err
	}
	eta, err := trial.SuggestLogUniform("eta", 1e-6, 1)
	if err != nil {
		return -1, err
	}
	latent, err := trial.SuggestInt("latent", 1, 16)
	if err != nil {
		return -1, err
	}
	number, err := trial.Number()
	if err != nil {
		return -1, err
	}
	jsonMetaPath := path.Join(workPath, fmt.Sprintf("ffm-meta-%d.json", number))

	ctx := trial.GetContext()
	cmd := exec.CommandContext(
		ctx,
		libffmPath,
		"-p", testDataPath,
		"--auto-stop", "--auto-stop-threshold", "3",
		"-l", fmt.Sprintf("%f", lmd),
		"-r", fmt.Sprintf("%f", eta),
		"-k", fmt.Sprintf("%d", latent),
		"-t", "500",
		"--json-meta", jsonMetaPath,
		dataPath,
	)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	_ = cmd.Run() // ignore because ffm-train exited with 1 when enabling early stopping.

	var result struct {
		BestIteration int     `json:"best_iteration"`
		BestVALoss    float64 `json:"best_va_loss"`
	}

	jsonStr, err := ioutil.ReadFile(jsonMetaPath)
	if err != nil {
		return -1, fmt.Errorf("failed to read json: %s", err)
	}
	err = json.Unmarshal(jsonStr, &result)
	if err != nil {
		return -1, fmt.Errorf("failed to read json: %s", err)
	}
	if result.BestIteration == 0 && result.BestVALoss == 0 {
		return -1, errors.New("failed to open json meta")
	}

	_ = trial.SetUserAttr("best_iteration", fmt.Sprintf("%d", result.BestIteration))
	_ = trial.SetUserAttr("stdout", stdout.String())
	_ = trial.SetUserAttr("stderr", stderr.String())
	return result.BestVALoss, nil
}
