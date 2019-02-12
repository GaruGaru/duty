package storage

import (
	"fmt"
	"github.com/GaruGaru/duty/task"
	bolt "go.etcd.io/bbolt"
)

type BoltDB struct {
	DB *bolt.DB
}

const Bucket = "duty-tasks"

func NewBoltDBStorage(path string) (*BoltDB, error) {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}
	return &BoltDB{DB: db}, err
}

func (s BoltDB) Store(task task.ScheduledTask) error {
	return s.DB.Update(func(tx *bolt.Tx) error {

		bucket, err := tx.CreateBucketIfNotExists([]byte(Bucket))

		if err != nil {
			return err
		}

		serialized, err := Serialize(task)

		if err != nil {
			return err
		}

		return bucket.Put([]byte(task.ID), serialized)
	})
}

func (s BoltDB) Status(id string) (task.ScheduledTask, error) {
	var task task.ScheduledTask

	err := s.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Bucket))

		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}

		payload := bucket.Get([]byte(id))

		if payload == nil || len(payload) == 0 {
			return fmt.Errorf("no task found with id " + id)
		}

		deserialized, err := Deserialize(payload)
		if err != nil {
			return err

		}
		task = deserialized

		return nil
	})

	return task, err
}

func (s BoltDB) ListByType(types string) ([]task.ScheduledTask, error) {
	tasks, err := s.ListAll()

	filtered := make([]task.ScheduledTask, 0)

	if err != nil {
		return nil, err
	}

	for _, task := range tasks {
		filtered = append(filtered, task)
	}

	return filtered, nil
}

func (s BoltDB) ListAll() ([]task.ScheduledTask, error) {
	tasks := make([]task.ScheduledTask, 0)

	err := s.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Bucket))

		if bucket == nil {
			return nil
		}

		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			scheduledTask, err := Deserialize(v)

			if err != nil {
				return err
			}

			tasks = append(tasks, scheduledTask)

		}

		return nil
	})

	return tasks, err
}

func (s BoltDB) Delete(id string) (bool, error) {
	deleted := false
	err := s.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Bucket))

		if bucket == nil {
			return nil
		}

		err := bucket.Delete([]byte(id))
		if err != nil {
			return err
		}

		deleted = true
		return nil
	})
	return deleted, err
}

func (s BoltDB) Update(task task.ScheduledTask, status task.Status) error {
	task, err := s.Status(task.ID)

	if err != nil {
		return err
	}

	task.Status = status

	return s.Store(task)
}

func (s BoltDB) Close() {
	_ = s.DB.Close()
}
