package dbbolt

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

func GetDB() *bolt.DB {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		return nil
	}

	return db
}

func main() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		bucketDb, err := tx.CreateBucketIfNotExists([]byte("evo_erp"))
		if err != nil {
			return err
		}
		bucketQ, err := bucketDb.CreateBucketIfNotExists([]byte("erp_uh"))
		if err != nil {
			return err
		}

		bucketN1, err := bucketQ.CreateBucketIfNotExists([]byte("n1"))
		if err != nil {
			return err
		}
		err = bucketN1.Put([]byte("msg1"), []byte("1"))
		if err != nil {
			return err
		}
		err = bucketN1.Put([]byte("msg32"), []byte("2"))
		if err != nil {
			return err
		}
		err = bucketN1.Put([]byte("msg3"), []byte("3"))
		if err != nil {
			return err
		}

		bucketN2, err := bucketQ.CreateBucketIfNotExists([]byte("n2"))
		if err != nil {
			return err
		}
		err = bucketN2.Put([]byte("msg1"), []byte("11"))
		if err != nil {
			return err
		}
		err = bucketN2.Put([]byte("msg22"), []byte("22"))
		if err != nil {
			return err
		}
		err = bucketN2.Put([]byte("msg3"), []byte("33"))
		if err != nil {
			return err
		}

		return nil
	})

	db.Update(func(tx *bolt.Tx) error {
		bucketDb := tx.Bucket([]byte("evo_erp"))
		if bucketDb == nil {
			return fmt.Errorf("bucket not found")
		}

		bucketQ := bucketDb.Bucket([]byte("erp_uh"))
		if bucketQ == nil {
			return fmt.Errorf("bucket not found")
		}

		bucketQ.ForEach(func(k, v []byte) error {

			upl := bucketQ.Bucket(k)
			fmt.Println(string(k))
			//upl.ForEach(func(k, v []byte) error {
			//	fmt.Println(string(k), string(v))
			//	upl.Delete(k)
			//
			//	return nil
			//})
			c := upl.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				upl.Delete(k)
				fmt.Println(string(k), string(v))
			}

			return nil
		})

		bucketQ.ForEach(func(k, v []byte) error {

			upl := bucketQ.Bucket(k)
			fmt.Println(string(k))
			c := upl.Cursor()
			fk, _ := c.First()

			if fk == nil {
				bucketQ.DeleteBucket(k)
			}
			return nil
		})

		return nil
	})

	db.View(func(tx *bolt.Tx) error {
		bucketDb := tx.Bucket([]byte("evo_erp"))
		if bucketDb == nil {
			return fmt.Errorf("bucket not found")
		}

		bucketQ := bucketDb.Bucket([]byte("erp_uh"))
		if bucketQ == nil {
			return fmt.Errorf("bucket not found")
		}

		c := bucketQ.Cursor()
		fk, _ := c.First()
		if fk == nil {
			return nil
		}
		return nil
	})
}
