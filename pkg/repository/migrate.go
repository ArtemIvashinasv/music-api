package repository

import (
    "database/sql"
    "fmt"
    "io/ioutil"
    "log"
    "path/filepath"
)

func Migrate(db *sql.DB) error {
    // Начинаем транзакцию
    tx, err := db.Begin()
    if err != nil {
        return fmt.Errorf("ошибка начала транзакции: %w", err)
    }

    // Обработчик для отката транзакции в случае ошибки
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

    // Читаем запрос из файла для создания таблицы миграций
    migrationsTableFilePath := filepath.Join("pkg/repository", "create_migrations_table.sql")
    migrationsTableQuery, err := ioutil.ReadFile(migrationsTableFilePath)
    if err != nil {
        return fmt.Errorf("ошибка чтения файла миграции: %w", err)
    }

    // Создаем таблицу миграций, если она не существует
    _, err = tx.Exec(string(migrationsTableQuery))
    if err != nil {
        return fmt.Errorf("ошибка создания таблицы миграций: %w", err)
    }

    // Получаем текущую версию миграции
    var currentVersion int
    err = tx.QueryRow("SELECT COALESCE(MAX(version), 0) FROM migrations").Scan(&currentVersion)
    if err != nil {
        return fmt.Errorf("ошибка получения текущей версии миграции: %w", err)
    }

    // Определяем целевую версию миграции
    const targetVersion = 1

    if currentVersion < targetVersion {
        // Читаем запрос из файла для создания таблицы песен
        migrationFilePath := filepath.Join("migrations", "create_songs_table.sql")
        migrationQuery, err := ioutil.ReadFile(migrationFilePath)
        if err != nil {
            return fmt.Errorf("ошибка чтения файла миграции: %w", err)
        }

        // Выполняем миграцию
        _, err = tx.Exec(string(migrationQuery))
        if err != nil {
            return fmt.Errorf("ошибка выполнения миграции: %w", err)
        }

        // Обновляем версию миграции
        _, err = tx.Exec("INSERT INTO migrations (version) VALUES ($1)", targetVersion)
        if err != nil {
            return fmt.Errorf("ошибка обновления версии миграции: %w", err)
        }

        log.Printf("Миграция успешно применена: версия %d", targetVersion)
    }

    // Фиксируем транзакцию
    err = tx.Commit()
    if err != nil {
        return fmt.Errorf("ошибка фиксации транзакции: %w", err)
    }

    return nil
}
