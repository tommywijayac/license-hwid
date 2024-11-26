package repo

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/google/uuid"
)

type Repo interface {
	Add(log LicenseLog)
}

type repo struct {
	mutex sync.Mutex
	db    Database
}

const dbPath string = "./repo/license_log.json"

type Database struct {
	Data []LicenseLog `json:"data"`
}

type LicenseLog struct {
	ID              uuid.UUID `json:"id"`
	HashedMachineID []byte    `json:"hashed_machine_id"`
	LicenseFilepath string    `json:"license_filepath"`
	HardwareLabel   string    `json:"hardware_label"`
	CreatedTime     string    `json:"created_time"`
}

func New() Repo {
	f, err := os.Open(dbPath)
	if err != nil {
		panic(err)
	}

	db := Database{}
	err = json.NewDecoder(f).Decode(&db)
	if err != nil {
		panic(err)
	}

	return &repo{
		db: db,
	}
}

func (r *repo) Add(log LicenseLog) {
	// serialize write
	r.mutex.Lock()
	defer r.mutex.Unlock()

	f, err := os.OpenFile(dbPath, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	log.ID = uuid.New()
	r.db.Data = append(r.db.Data, log)

	jsonb, err := json.MarshalIndent(r.db, "", "\t")
	if err != nil {
		panic(err)
	}

	_, err = f.Write(jsonb)
	if err != nil {
		panic(err)
	}
}
