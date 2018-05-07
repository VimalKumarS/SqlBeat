package beater

import (
	"database/sql"
	"fmt"
	"time"

	"sqlbeat/config"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

type Sqlbeat struct {
	done          chan struct{}
	config        config.Config
	client        beat.Client
	lastIndexTime time.Time
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Sqlbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

// Run method
func (bt *Sqlbeat) Run(b *beat.Beat) error {
	logp.Info("sqlbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	for {
		now := time.Now()
		err := bt.beat(b)
		if err != nil {
			return err
		}
		bt.lastIndexTime = now // use last index time to query latest value
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		// event := beat.Event{
		// 	Timestamp: time.Now(),
		// 	Fields: common.MapStr{
		// 		"type":    b.Info.Name,
		// 		"counter": counter,
		// 	},
		// }
		// bt.client.Publish(event)
		logp.Info("Event sent")
		// counter++
	}
}

// Stop method
func (bt *Sqlbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

func (bt *Sqlbeat) beat(b *beat.Beat) error {
	connString := ""
	connString = fmt.Sprintf("server=%v;user id=%v;password=%v;database=%v",
		bt.config.Hostname, bt.config.Username, bt.config.Password, bt.config.Database)

	db, err := sql.Open("mssql", connString)
	if err != nil {
		logp.Err(err.Error())
		return err
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		logp.Err(err.Error())
		fmt.Println("Cannot connect: ", err.Error())
		return err
	}

	dtNow := time.Now()
	fmt.Println(bt.lastIndexTime.Format("2006-01-02"))
	rows, err := db.Query(bt.config.Queries)
	defer rows.Close()
	if err != nil {
		logp.Err(err.Error())
		return err
	}

	// Populate columns array
	cols, err := rows.Columns()
	if err != nil {
		logp.Err(err.Error())
		return err
	}

	if cols == nil {
		return nil
	}
	vals := make([]sql.RawBytes, len(cols))
	scanArgs := make([]interface{}, len(vals))
	fmt.Println(dtNow)
	for i := range vals {
		scanArgs[i] = &vals[i]
	}

	for rows.Next() {

		mapstr := common.MapStr{
			"type": "category",
		}
		err = rows.Scan(scanArgs...)
		if err != nil {
			logp.Err(err.Error())
			fmt.Println(err)
			continue
		}
		id := ""
		for i, col := range vals {
			strColName := string(cols[i])
			strColValue := string(col)
			if strColName == "id" {
				id = strColValue
			}

			fmt.Println(strColName, strColValue)
			mapstr[strColName] = strColValue
		}
		event := beat.Event{
			Timestamp: time.Now(),
			Fields:    mapstr,
			Meta: common.MapStr{
				"id": id,
			},
			Private: nil,
		}
		event.SetID(id)
		bt.client.Publish(event) // publish event to elastic search
	}
	if rows.Err() != nil {
		logp.Err(rows.Err().Error())
		return rows.Err()
	}
	// Great success!
	return nil
}
