package paxos

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/fudute/GoPaxos/sm"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var ErrorBadLogFormat = errors.New("bad log format")

func InitDB() {
	DB = NewLevelDB("../db")
	DB.Restore(GetProposerInstance(), GetAcceptorInstance(), sm.GetKVStatMachineInstance())
}

// [key = index, set key value ]
// [key = index, delete key ]
// [key = index, nop ]

type LevelDB struct {
	db  *leveldb.DB
	opt *opt.Options
	ro  *opt.ReadOptions
	wo  *opt.WriteOptions
	cmp MyComparator
}

func NewLevelDB(path string) *LevelDB {

	db := &LevelDB{}
	db.opt = &opt.Options{}
	db.cmp = MyComparator{}
	db.opt.Comparer = db.cmp
	db.ro = &opt.ReadOptions{}
	db.wo = &opt.WriteOptions{}

	var err error
	db.db, err = leveldb.OpenFile(path, db.opt)

	if err != nil {
		log.Fatal(err)
	}
	return db
}

func (level *LevelDB) WriteLog(index int, entry *LogEntry) error {

	str := entry.String()
	key := strconv.Itoa(index)
	return level.db.Put([]byte(key), []byte(str), level.wo)
}

func (level *LevelDB) ReadLog(index int) (*LogEntry, error) {

	key := strconv.Itoa(index)

	value, err := level.db.Get([]byte(key), level.ro)

	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, ErrorNotFound
		}
		return nil, err
	}

	le, err := parseLog(string(value))

	if err != nil {
		return nil, err
	}

	return le, nil
}

func (level *LevelDB) Restore(p *Proposer, a *Acceptor, statMachine sm.StatMachine) error {

	iter := level.db.NewIterator(nil, nil)

	log.Println("start restore...")
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		le, err := parseLog(string(value))

		log.Println("restore: log entry [", key, "]  ", *le)

		if err != nil {
			return err
		}

		index, err := strconv.Atoi(string(key))
		if err != nil {
			log.Fatal("db format is incorrect")
		}
		if p.LogIndex == index && le.IsCommited {
			cmd := sm.Command{
				Index: p.LogIndex,
				Cmd:   le.AcceptedValue,
			}
			statMachine.Execute(cmd)
			p.LogIndex++
		}
	}
	return nil
}

func (level *LevelDB) Close() error {
	return level.db.Close()
}

func parseLog(str string) (*LogEntry, error) {
	var err error

	le := &LogEntry{}
	tokens := strings.Split(str, ":")

	if len(tokens) < 3 {
		return nil, ErrorBadLogFormat
	}
	le.MinProposal, err = strconv.Atoi(tokens[0])
	if err != nil {
		return nil, ErrorBadLogFormat
	}

	le.AcceptedProposal, err = strconv.Atoi(tokens[1])
	if err != nil {
		return nil, ErrorBadLogFormat
	}

	le.AcceptedValue = tokens[2]

	le.IsCommited, err = strconv.ParseBool(tokens[3])
	if err != nil {
		return nil, ErrorBadLogFormat
	}

	return le, nil
}

func extractCommand(value string) string {
	var ret string
	first := -2
	for i := 0; i < len(value); i++ {
		if value[i] == ':' {
			if first == -2 {
				first = -1
			} else if first == -1 {
				first = i
			} else {
				ret = value[first+1 : i]
			}
		}
	}

	return ret
}

func (db *LevelDB) PrintLog(fileName string) {
	view, err := db.db.GetSnapshot()
	if err != nil {
		log.Println("Get snap shot error:", err)
	}

	iter := view.NewIterator(nil, nil)
	file := openNewEmptyFile("/home/log/" + fileName)
	defer file.Close()

	log_index := 0
	for iter.Next() {
		key, value := iter.Key(), iter.Value()
		le, err := parseLog(string(value))
		if err != nil {
			log.Printf("parse log failed, key :[%v] ,value :%v\n", key, value)
			continue
		}
		// TODO
		// > 1974 17:8:SET ntup xoyhdzdnxe:true
		// ---
		// > 1974 17:10:SET ntup xoyhdzdnxe:true
		command := extractCommand(string(value))
		db_index, err := strconv.Atoi(string(key))
		if err != nil {
			log.Printf("log index [%v] cant convert to int\n", string(key))
		}
		for i := log_index; i < db_index; i++ {
			_, err = file.WriteString(strconv.Itoa(i) + " " + strconv.FormatBool(false) + " UNKNOWN\n")
			if err != nil {
				log.Println("write log error:", err)
			}
		}

		_, err = file.WriteString(string(key) + " " + strconv.FormatBool(le.IsCommited) + " " + command + "\n")

		if err != nil {
			log.Printf("cant write commited log index:[%v] to file\n", string(key))
			continue
		}
		log_index++

	}
	fmt.Printf("print all logs(total: %v) to ~/paxos_logs/n%v/%v"+
		" successed, please check it later\n", log_index, GetServerID(), fileName)
}
