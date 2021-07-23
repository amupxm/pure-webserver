package database

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/amupxm/pure-webserver/config"
)

type (
	Database interface {
		// WriteToCollection writes a collection to the database
		WriteToCollection(collection interface{}) error
		// GetFromCollection returns a collection from the database
		GetFromCollection(collection interface{}) (DBInnerModel, error)
		// UpdateCollection updates a collection in the database
		UpdateCollection(collectionData interface{}) error
		// getCollections returns a collections of DbModel (creates one if does not exist)
		getCollections(collection interface{}) (*DbModelCollection, error)
		// getCollectionName returns collection name as string
		getCollectionName(collection interface{}) string
		// readDatabase reads the database with io
		readDatabase() (*DbModelCollection, error)
		// writeDatabase writes the database with io
		writeDatabase(dbCollection *DbModelCollection) error
	}
	database struct {
		// The lock for the database (experimental to provide acid)
		lock sync.RWMutex
	}
	DbModel struct {
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at"`
		Deleted   bool       `json:"deleted"`
		Id        string     `json:"id"`
	}
	DbModelCollection struct {
		Items map[string]DBInnerModel `json:"items"`
		Meta  struct {
			Total int `json:"total"`
		} `json:"meta"`
		DataIndexes map[string]int `json:"data_indexes"` // to save count of items stored in collection
	}
	DbModelCollectionInterface interface {
		Where(fieldName string, value interface{}) *DBInnerModel
		All() *DBInnerModel
		Update(data interface{}) (*DBInnerModel, error)
	}
	DBInnerModel []interface {
	}
)

// NewDatabase creates a new Database instance
func NewDatabase() Database {
	return &database{
		lock: sync.RWMutex{},
	}
}

// readDatabase reads the database
func (db *database) readDatabase() (*DbModelCollection, error) {
	var result DbModelCollection
	b, err := ioutil.ReadFile(config.AppConf.DatabaseConfig.BucketName)
	if err != nil {
		return &result, err
	}
	err = json.Unmarshal(b, &result)
	if err != nil {
		result = DbModelCollection{Items: make(map[string]DBInnerModel), DataIndexes: make(map[string]int)}
	}
	return &result, nil
}

// writeDatabase writes the database
func (db *database) writeDatabase(dbCollection *DbModelCollection) error {
	marshaledDbData, err := json.Marshal(dbCollection)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(config.AppConf.DatabaseConfig.BucketName, []byte(marshaledDbData), 0666)
	if err != nil {
		return err
	}
	return nil
}

// getCollection returns a collections of DbModel (creates one if does not exist)
func (db *database) getCollections(collection interface{}) (*DbModelCollection, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()
	dbCollection, err := db.readDatabase()
	if err != nil {
		return nil, err
	}
	// create dbCollection if not exists
	collectionName := db.getCollectionName(collection)
	// check collectionName exists in dbCollection

	if _, ok := dbCollection.Items[collectionName]; !ok {

		// if not exists, create a new collection
		if dbCollection.Items == nil {
			dbCollection.Items = make(map[string]DBInnerModel)
		}
		if dbCollection.DataIndexes == nil {
			dbCollection.DataIndexes = make(map[string]int)
		}
		// create collection
		dbCollection.Items[collectionName] = []interface{}{}
		dbCollection.DataIndexes[collectionName] = 0
	}
	return dbCollection, nil
}

// getCollectionName  returns collection name as string
func (db *database) getCollectionName(collection interface{}) string {
	return reflect.TypeOf(collection).String()
}

// WriteToCollection writes a collection to the database
func (db *database) WriteToCollection(collection interface{}) error {
	dbCollection, err := db.getCollections(collection)
	if err != nil {
		return err
	}

	lastId := dbCollection.DataIndexes[db.getCollectionName(collection)]
	f := reflect.Indirect(reflect.ValueOf(collection)).FieldByName("DbModel")
	f.FieldByName("Id").SetString(strconv.Itoa(lastId + 1))
	f.FieldByName("CreatedAt").Set(reflect.ValueOf(time.Now()))
	f.FieldByName("UpdatedAt").Set(reflect.ValueOf(time.Now()))

	jsonC, err := json.Marshal(f.Interface())
	if err != nil {
		return err
	}
	_ = json.Unmarshal(jsonC, collection)

	c := dbCollection.Items[db.getCollectionName(collection)]
	c = append(c, collection)
	dbCollection.Items[db.getCollectionName(collection)] = c
	dbCollection.DataIndexes[db.getCollectionName(collection)]++
	return db.writeDatabase(dbCollection)
}

// UpdateCollection updates a collection in the database
func (db *database) UpdateCollection(typecollectionData interface{}) error {
	reflection := reflect.TypeOf(typecollectionData)
	var typeTarget interface{}
	if reflection.Kind() != reflect.Array && reflection.Kind() == reflect.Slice {
		typeTarget = typecollectionData
	} else {
		// get typeOf firstIndex
		typeTarget = reflect.Indirect(reflect.ValueOf(typecollectionData)).Index(0).Interface()
	}
	dbCollection, err := db.getCollections(typeTarget)
	if err != nil {
		return err
	}
	tmpName := db.getCollectionName(typeTarget)
	if !strings.HasPrefix(tmpName, "*") {
		tmpName = "*" + tmpName
	}
	c := dbCollection.Items[db.getCollectionName(tmpName)]
	// if typecollectionData is array, then update all items
	if reflection.Kind() == reflect.Array || reflection.Kind() == reflect.Slice || reflection.Kind() == reflect.Ptr {
		for i := 0; i < reflect.Indirect(reflect.ValueOf(typecollectionData)).Len(); i++ {
			f := reflect.Indirect(reflect.ValueOf(typecollectionData)).Index(i)
			f.FieldByName("UpdatedAt").Set(reflect.ValueOf(time.Now()))
			jsonC, err := json.Marshal(f.Interface())
			if err != nil {
				return err
			}
			_ = json.Unmarshal(jsonC, f.Interface())
			c = append(c, f.Interface())
		}
		dbCollection.Items[tmpName] = c
		return db.writeDatabase(dbCollection)
	}
	return nil
}

// GetFromCollection returns a collection from the database
func (db *database) GetFromCollection(collection interface{}) (DBInnerModel, error) {
	dbCollection, err := db.getCollections(collection)
	if err != nil {
		return nil, err
	}
	collectionName := db.getCollectionName(collection)
	//TODO set kind of collection to DBInnerModel
	c := dbCollection.Items[collectionName]
	if len(c) == 0 {
		return nil, nil
	}

	return c, nil
}

// Where return one or more items from the database where the filter is true
func (dbm *DBInnerModel) Where(fieldName string, value interface{}) *DBInnerModel {
	var temp []interface{}

	valueKind := reflect.TypeOf(value).Kind()
	for _, v := range *dbm {
		// check v have field named fieldName
		if reflect.TypeOf(v).Kind() == reflect.Map {
			c := v.(map[string]interface{})[fieldName]
			field := reflect.ValueOf(c)
			if field.IsValid() {
				if field.Kind() == valueKind {
					switch valueKind {
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						if field.Int() == value.(int64) {
							temp = append(temp, v)
						}
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						if field.Uint() == value.(uint64) {
							temp = append(temp, v)
						}
					case reflect.Float32, reflect.Float64:
						if field.Float() == value.(float64) {
							temp = append(temp, v)
						}
					case reflect.String:
						if field.String() == value.(string) {
							temp = append(temp, v)
						}
					case reflect.Bool:
						if field.Bool() == value.(bool) {
							temp = append(temp, v)
						}
						// TODO : slice and time and ...
					}
				}
			}
		}

	}
	result := DBInnerModel{}
	result = append(result, temp...)

	return &result
}

// All return all items from the database
func (dbm *DBInnerModel) All() *DBInnerModel {
	return dbm
}
