package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func multiplyMatrices(a, b [][]int) [][]int {
	rowsA := len(a)
	colsA := len(a[0])
	colsB := len(b[0])

	result := make([][]int, rowsA)
	for i := range result {
		result[i] = make([]int, colsB)
	}

	for i := 0; i < rowsA; i++ {
		for j := 0; j < colsB; j++ {
			for k := 0; k < colsA; k++ {
				result[i][j] += a[i][k] * b[k][j]
			}
		}
	}

	return result
}

func multiplyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		matrixStrA := r.FormValue("matrixA")
		matrixStrB := r.FormValue("matrixB")

		matrixA := parseMatrix(matrixStrA)
		matrixB := parseMatrix(matrixStrB)

		if matrixA == nil || matrixB == nil {
			http.Error(w, "Invalid matrix format", http.StatusBadRequest)
			return
		}

		if len(matrixA[0]) != len(matrixB) {
			http.Error(w, "Matrices cannot be multiplied: invalid dimensions", http.StatusBadRequest)
			return
		}

		result := multiplyMatrices(matrixA, matrixB)

		outputHTML := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Matrix Multiplication</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f4f4f4;
					margin: 0;
					padding: 20px;
				}
				h1 {
					text-align: center;
				}
				form {
					width: 50%;
					margin: auto;
					background-color: #fff;
					padding: 20px;
					border-radius: 8px;
					box-shadow: 0 0 10px rgba(0,0,0,0.1);
				}
				textarea {
					width: 100%;
					height: 100px;
					resize: none;
					border-radius: 4px;
					padding: 5px;
					margin-bottom: 10px;
				}
				input[type="submit"] {
					padding: 8px 16px;
					font-size: 16px;
					border: none;
					border-radius: 4px;
					background-color: #4CAF50;
					color: white;
					cursor: pointer;
					margin-top: 10px;
				}
				pre {
					white-space: pre-wrap;
				}
			</style>
		</head>
		<body>
			<h1>Matrix Multiplication Result</h1>
			<div style="text-align: center;">
				<pre>`
		for _, row := range result {
			outputHTML += fmt.Sprintf("%v\n", row)
		}
		outputHTML += `
				</pre>
			</div>
		</body>
		</html>`
		fmt.Fprintf(w, outputHTML)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func parseMatrix(matrixStr string) [][]int {
	rowsStr := strings.Split(strings.TrimSpace(matrixStr), "\n")
	rows := make([][]int, len(rowsStr))

	for i, rowStr := range rowsStr {
		colStr := strings.Fields(rowStr)
		row := make([]int, len(colStr))

		for j, valStr := range colStr {
			val, err := strconv.Atoi(valStr)
			if err != nil {
				return nil
			}
			row[j] = val
		}
		rows[i] = row
	}

	return rows
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	homeTemplate := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Matrix Multiplication</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #f4f4f4;
				margin: 0;
				padding: 20px;
			}
			h1 {
				text-align: center;
			}
			form {
				width: 50%;
				margin: auto;
				background-color: #fff;
				padding: 20px;
				border-radius: 8px;
				box-shadow: 0 0 10px rgba(0,0,0,0.1);
			}
			textarea {
				width: 100%;
				height: 100px;
				resize: none;
				border-radius: 4px;
				padding: 5px;
				margin-bottom: 10px;
			}
			input[type="submit"] {
				padding: 8px 16px;
				font-size: 16px;
				border: none;
				border-radius: 4px;
				background-color: #4CAF50;
				color: white;
				cursor: pointer;
				margin-top: 10px;
			}
		</style>
	</head>
	<body>
		<h1>Matrix Multiplication</h1>
		<form action="/multiply" method="post">
			<h2>Enter Matrix A:</h2>
			<textarea name="matrixA" rows="5" cols="40"></textarea><br/><br/>
			<h2>Enter Matrix B:</h2>
			<textarea name="matrixB" rows="5" cols="40"></textarea><br/><br/>
			<input type="submit" value="Умножить">
		</form>
	</body>
	</html>`
	fmt.Fprintf(w, homeTemplate)
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/multiply", multiplyHandler)

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
