package main

import (
    "fmt"

    "github.com/google/uuid"
    "gorm.io/driver/clickhouse"
    "gorm.io/gorm"
)

type User struct {
    ID   uuid.UUID `gorm:"type:UUID;default:generateUUIDv4()"`
    Name string
    Age  int
}

func main() {
    dsn := "clickhouse://gorm:gorm@localhost:9942/gorm?dial_timeout=10s&read_timeout=20s"
    db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    // Recreate table
    db.Migrator().DropTable(&User{})
    db.AutoMigrate(&User{})

    // Insert
    db.Create(&User{Name: "Angeliz", Age: 18})

    // Select by name
    var users []User
    db.Find(&users, "name = ?", "Angeliz")
    fmt.Println("Users found:", len(users))

    // Select by UUID (example)
    if len(users) > 0 {
        var u User
        db.Find(&u, "id = ?", users[0].ID)
        fmt.Println("Found by UUID:", u)
    }
}
