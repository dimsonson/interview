package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Intervals definition.
const (
	min1 int = 1
	min2 int = 2
	min5 int = 5
)

// price is stricture describes of input line.
type price struct {
	tiker     string
	price     int
	priceTime time.Time
}

// candle is stricture describes of output line.
type candle struct {
	tiker      string
	open       int
	max        int
	min        int
	close      int
	candleTime time.Time
	interval   int
}

// key is a structure of index key for data manipulation.
type key struct {
	tiker      string
	candleTime time.Time
}

// cndls is a sclice for use by Sort interface.
type cndls []candle

func main() {
	candles := make([]candle, 0)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	prices := make([]price, 0)

	var count int
	for scanner.Scan() {
		// solution
		scanTxt := scanner.Text()
		arr := strings.Split(scanTxt, ",")
		tm, err := time.Parse(time.RFC3339, arr[2])
		if err != nil {
			log.Fatal(err)
		}
		pr, err := strconv.ParseFloat(arr[1], 64)
		if err != nil {
			log.Fatal(err)
		}
		prices = append(prices, price{tiker: arr[0], price: int(pr * 100), priceTime: tm})
		count++
	}

	// solution
	intervals := []int{min1, min2, min5}

	for _, interval := range intervals {
		index := make(map[key][]int, 0)
		for _, price := range prices {
			k := key{tiker: price.tiker, candleTime: price.priceTime.Truncate(time.Duration(interval) * time.Minute)}
			_, ok := index[k]
			if !ok {
				index[k] = []int{price.price}
			}
			if ok {
				index[k] = append(index[k], price.price)
			}
		}
		for k, v := range index {
			max, min := maxMin(v)
			c := candle{
				tiker:      k.tiker,
				open:       v[0],
				max:        max,
				min:        min,
				close:      v[len(v)-1],
				candleTime: k.candleTime,
				interval:   interval,
			}
			candles = append(candles, c)
		}
	}

	By(sortorder).Sort(candles)

	w := csv.NewWriter(os.Stdout)
	w.Comma = ','
	defer w.Flush()

	for _, candle := range candles {
		if err := w.Write(candle.ToCSV()); err != nil {
			log.Fatal(err)
		}
	}
}

// ToCSV prepare candle data to CSV output.
func (c *candle) ToCSV() []string {
	candleCSV := make([]string, 4) // output fields: tiker, (open, high, low, close), candle time, candle interval
	candleCSV[0] = c.tiker
	candleCSV[1] = fmt.Sprintf("%2.2f,%2.2f,%2.2f,%2.2f", float64(c.open)/100, float64(c.max)/100, float64(c.min)/100, float64(c.close)/100) //strings.Join(c.open, ",") // open, high, low, close
	candleCSV[2] = c.candleTime.Format(time.RFC3339)
	candleCSV[3] = fmt.Sprint(c.interval, "min")
	return candleCSV
}

// maxMin provides max and min from slice.
func maxMin(arr []int) (int, int) {
	var max int
	var min int
	for idx, v := range arr {
		if idx == 0 {
			min = v
		}
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	return max, min
}

// sortorder is comparator.
func sortorder(cndl1, cndl2 *candle) bool {
	return cndl1.tiker < cndl2.tiker ||
		cndl1.tiker == cndl2.tiker && cndl1.interval < cndl2.interval ||
		cndl1.tiker == cndl2.tiker && cndl1.interval == cndl2.interval && cndl1.candleTime.Before(cndl2.candleTime)
}

// cndlsSorter implements Interface interface.
type cndlsSorter struct {
	cndls
	by func(cnd1, cnd2 *candle) bool
}

type By func(cnd1, cnd2 *candle) bool

// Sort construct and run.
func (cnd By) Sort(cndls []candle) {
	cndlsSorter := &cndlsSorter{
		cndls: cndls,
		by:    cnd,
	}
	sort.Sort(cndlsSorter)
}

// Len is part of sort.Interface.
func (cnd *cndlsSorter) Len() int {
	return len(cnd.cndls)
}

// Swap is part of sort.Interface.
func (cnd *cndlsSorter) Swap(i, j int) {
	cnd.cndls[i], cnd.cndls[j] = cnd.cndls[j], cnd.cndls[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (cnd *cndlsSorter) Less(i, j int) bool {
	return cnd.by(&cnd.cndls[i], &cnd.cndls[j])
}
