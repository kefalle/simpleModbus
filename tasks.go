package main

import (
	"fmt"
	"log"
	"math"
	"os"
)

func Tasks() {
	// Задача 1
	task1 := func(x float64, y float64) {
		var r float64 = 2

		res := math.Abs(x*x + y*y - r*r)
		if res <= 1/1000 {
			log.Printf("task1: (%f,%f) is success", x, y)
		} else {
			log.Printf("task1: (%f,%f) is success", x, y)
		}
	}

	task1(0, 2)
	task1(1.5, 0.7)
	task1(1, 1)
	task1(3, 0)

	task2 := func(a float64, b float64, c float64) float64 {
		return math.Max(math.Min(a, b), c)
	}

	log.Printf("task2: F=%f", task2(1, 2, 3))

	task3 := func(r float64, s float64) bool {
		radius := math.Sqrt(r / math.Pi)
		maxS := 2 * radius * radius
		res := s <= maxS

		log.Printf("task3: r=%f s=%f is %t", r, s, res)

		return res
	}

	task3(70, 36.74)
	task3(0.86, 0.74)

	task21 := func() {
		const n uint = 10

		var m, w uint
		print("Enter boys count:")
		fmt.Scanf("%d", &m)

		print("Enter girls count:")
		fmt.Scanf("%d", &w)

		if m+w != n {
			log.Printf("Girls and boys count must be %d", n)
			return
		}

		boys := make([]float64, m)
		for i := range boys {
			log.Printf("Enter height for boy %d:", i+1)
			fmt.Scanf("%f", &boys[i])
		}

		girls := make([]float64, w)
		for i := range girls {
			log.Printf("Enter height for girl %d:", i+1)
			fmt.Scanf("%f", &girls[i])
		}

		avg := func(arr []float64) float64 {
			var sum float64
			for _, b := range arr {
				sum += b
			}
			return sum / float64(len(arr))
		}

		fmt.Printf("Avg boys h=%f, avg girls h=%f", avg(boys), avg(girls))
	}

	task21()

	os.Exit(0)
}
