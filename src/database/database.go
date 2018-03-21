package database



connStr := "user=pqgotest dbname=pqgotest sslmode=verify-full"
db, err := sql.Open("postgres", connStr)
if err != nil {
log.Fatal(err)
}

age := 21
