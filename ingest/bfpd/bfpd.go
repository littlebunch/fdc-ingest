// Package bfpd implements an Ingest for Branded Food Products
package bfpd

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/littlebunch/fdc-api/ds"
	fdc "github.com/littlebunch/fdc-api/model"
	"github.com/littlebunch/fdc-ingest/ingest"
	"github.com/littlebunch/fdc-ingest/ingest/dictionaries"
)

var (
	cnts ingest.Counts
	err  error
	//gbucket string
)

// Bfpd for implementing the interface
type Bfpd struct {
	Doctype string
}
type line struct {
	id         int
	restOfLine string
}
type f struct {
	FdcID string `json:"fdcId" binding:"required"`
}

// ProcessFiles loads a set of Branded Food Products csv
func (p Bfpd) ProcessFiles(path string, dc ds.DataSource, bucket string) error {

	/*	var (
			dt   *fdc.DocType
			il   []interface{}
			food fdc.Food
			s    []fdc.Serving
			err  error
		)
		gbucket = bucket
		if il, err = dc.GetDictionary(gbucket, dt.ToString(fdc.FGGPC), 0, 500); err != nil {
			return err
		}
		fgrp := dictionaries.InitBrandedFoodGroupInfoMap(il)
		// read food metadata
		metadataChan := make(chan *line)
		go reader(path+"food.csv", metadataChan)
		// read branded foods details
		brandedChan := make(chan *line)
		go reader(path+"branded_food.csv", brandedChan)
		// join the two data streams
		mergedLinesChan := make(chan *line)
		go joiner(metadataChan, brandedChan, mergedLinesChan)
		// process the merge stream
		var buf bytes.Buffer
		r := csv.NewReader(&buf)
		fgid := 0

		for l := range mergedLinesChan {
			buf.WriteString(fmt.Sprintf("%v,%v", l.id, l.restOfLine))
			record, _ := r.Read()
			fgid++
			if fgid%10000 == 0 {
				log.Println(fgid)
			}

			if rc := dc.FoodExists(record[0]); rc {
				continue
			} else { // create a new food
				s = nil
				fmt.Printf("%v,%v\n", l.id, l.restOfLine)
				pubdate, err := time.Parse("2006-01-02", record[4])
				if err != nil {
					log.Println(err)
				}
				food.ID = record[0]
				food.FdcID = record[0]
				food.Description = record[2]
				food.PublicationDate = pubdate
				food.Manufacturer = record[5]
				food.Upc = record[6]
				food.Ingredients = record[7]
				cnts.Foods++
				if cnts.Foods%10000 == 0 {
					log.Println("Foods Count = ", cnts.Foods)
				}
				a, err := strconv.ParseFloat(record[8], 32)
				if err != nil {
					log.Println(record[0] + ": can't parse serving amount " + record[8])
				} else {
					s = append(s, fdc.Serving{
						Nutrientbasis: record[9],
						Description:   record[10],
						Servingamount: float32(a),
					})
					food.Servings = s
				}
				food.Source = record[12]
				if record[13] != "" {
					food.ModifiedDate, _ = time.Parse("2006-01-02", record[13])
				}
				if record[14] != "" {
					food.AvailableDate, _ = time.Parse("2006-01-02", record[14])
				}
				if record[16] != "" {
					food.DiscontinueDate, _ = time.Parse("2006-01-02", record[16])
				}
				food.Country = record[15]
				food.Type = dt.ToString(fdc.FOOD)
				if record[11] != "" {
					_, fg := fgrp[record[11]]
					if !fg {
						fgid++
						fgrp[record[11]] = fdc.FoodGroup{ID: int32(fgid), Description: record[11], Type: dt.ToString(fdc.FGGPC)}
					}
					food.Group = &fdc.FoodGroup{ID: int32(fgrp[record[11]].ID), Description: fgrp[record[11]].Description, Type: fgrp[record[11]].Type}
				} else {
					food.Group = nil
				}
				// first remove any existing versions for this GTIN/UPC code
				removeVersions(food.Upc, bucket, dc)
				if err = dc.Update(record[0], food); err != nil {
					log.Printf("Update %s failed: %v", record[0], err)
				}

			}

		}
	*/
	if err = nutrients(path, bucket, dc); err != nil {
		fmt.Printf("nutrient load failed: %v", err)
	}

	log.Printf("Finished.  Counts: %d Foods %d Servings %d Nutrients\n", cnts.Foods, cnts.Servings, cnts.Nutrients)
	return err
}
func removeVersions(upc string, bucket string, dc ds.DataSource) {

	var (
		r   []interface{}
		fid f
		j   []byte
	)

	q := fmt.Sprintf("SELECT fdcId from %s where upc = \"%s\" AND type=\"FOOD\"", bucket, upc)
	if err := dc.Query(q, &r); err != nil {
		log.Printf("%v\n", err)
		return
	}
	for i := range r {
		if j, err = json.Marshal(r[i]); err != nil {
			log.Printf("%s %v %v\n", upc, j, err)
		}
		if err = json.Unmarshal(j, &fid); err != nil {
			log.Printf("%s %s %v\n", upc, string(j), err)
		}
		log.Printf("Removed %s\n", fid.FdcID)
		if err = dc.Remove(fid.FdcID); err != nil {
			log.Printf("Cannot remove %s\n", fid.FdcID)
		}
	}
	return

}

func nutrients(path string, gbucket string, dc ds.DataSource) error {
	var (
		dt          *fdc.DocType
		food        fdc.Food
		cid, source string
	)
	fn := path + "food_nutrient.csv"
	f, err := os.Open(fn)
	if err != nil {
		return err
	}

	r := csv.NewReader(f)
	var (
		n  []fdc.NutrientData
		il []interface{}
	)
	q := fmt.Sprintf("select gd.* from %s as gd where type='%s' offset %d limit %d", gbucket, dt.ToString(fdc.NUT), 0, 500)
	fmt.Println(q)
	if il, err = dc.GetDictionary(gbucket, dt.ToString(fdc.NUT), 0, 500); err != nil {
		return err
	}

	nutmap := dictionaries.InitNutrientInfoMap(il)

	if il, err = dc.GetDictionary(gbucket, dt.ToString(fdc.DERV), 0, 500); err != nil {
		return err
	}
	dlmap := dictionaries.InitDerivationInfoMap(il)
	processit := true
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		id := record[1]
		v, err := strconv.ParseInt(record[2], 0, 32)
		if err != nil {
			log.Println(record[0] + ": can't parse nutrient no " + record[1])
		}
		if processit = dc.FoodExists(id); !processit {
			// delete this record if the parent food doesn't exist
			nid := fmt.Sprintf("%s_%d", id, nutmap[uint(v)].Nutrientno)
			if err = dc.Remove(nid); err != nil {
				log.Printf("Problem with removing %s: %v\n", nid, err)
				continue
			}
		}
		cnts.Nutrients++
		w, err := strconv.ParseFloat(record[3], 32)
		if err != nil {
			log.Println(record[0] + ": can't parse value " + record[4])
		}

		d, err := strconv.ParseInt(record[5], 0, 32)
		if err != nil {
			log.Println(record[5] + ": can't parse derivation no " + record[1])
		}
		var dv *fdc.Derivation
		if dlmap[uint(d)].Code != "" {
			dv = &fdc.Derivation{ID: dlmap[uint(d)].ID, Code: dlmap[uint(d)].Code, Type: dt.ToString(fdc.DERV), Description: dlmap[uint(d)].Description}
		} else {
			dv = nil
		}
		if cid != id {
			if err = dc.Get(id, &food); err != nil {
				log.Printf("Cannot find %s %v", id, err)
			}
			cid = id
			source = food.Source

		}

		n = append(n, fdc.NutrientData{
			ID:         fmt.Sprintf("%s_%d", id, nutmap[uint(v)].Nutrientno),
			FdcID:      id,
			Nutrientno: nutmap[uint(v)].Nutrientno,
			Value:      float32(w),
			Nutrient:   nutmap[uint(v)].Name,
			Unit:       nutmap[uint(v)].Unit,
			Derivation: dv,
			Type:       dt.ToString(fdc.NUTDATA),
			Source:     source,
		})
		if cnts.Nutrients%1000 == 0 {
			log.Println("Nutrients Count = ", cnts.Nutrients)
			err := dc.Bulk(&n)
			if err != nil {
				log.Printf("Bulk insert failed: %v\n", err)
			}
			n = nil
		}

	}

	return nil
}
func reader(fname string, out chan<- *line) {
	defer close(out) // close channel on return

	// open the file
	file, err := os.Open(fname)
	if err != nil {
		log.Fatalf("open: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	header := true
	for scanner.Scan() {
		var l line
		columns := strings.SplitN(scanner.Text(), ",", 2)
		// ignore first line (header)
		if header {
			header = false
			continue
		}
		// convert ID to integer for easier comparison
		id, err := strconv.Atoi(strings.ReplaceAll(columns[0], "\"", ""))
		if err != nil {
			log.Fatalf("ParseInt: %v", err)
		}
		l.id = id
		l.restOfLine = columns[1]
		// send the line to the channel
		out <- &l
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
func joiner(metadata, setIDs <-chan *line, out chan<- *line) {
	defer close(out) // close channel on return

	bf := &line{}
	for md := range metadata {
		sep := ","
		// add matching branded_foods.csv line (if left over from previous iteration)
		if bf.id == md.id {
			md.restOfLine += sep + bf.restOfLine
			sep = " "
		}
		// look for matching branded foods
		for bf = range setIDs {
			// add all branded_foods.csv with matching IDs
			if bf.id == md.id {
				md.restOfLine += sep + bf.restOfLine
				sep = " "
			} else if bf.id > md.id {
				break
			}
		}
		// send the augmented line into the channel
		out <- md
	}
}
