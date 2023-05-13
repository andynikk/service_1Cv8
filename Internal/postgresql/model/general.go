package model

// KeyContext тип для создания ключа контекста
type KeyContext string

// Updater интерфес для работы с объектами БД, отвечающих условиям контракта.
// На текущий момент это User, BankCard, BinaryData, PairLoginPassword, TextData
type Updater interface {
}

// Appender мапа для хранения объектов Updater
type Appender map[string]Updater
