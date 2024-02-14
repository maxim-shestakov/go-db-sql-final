package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel) //adding new parcel in db
	require.NoError(t, err)      //check
	require.NotEmpty(t, id)
	parcel.Number = id //number assignment after adding the database entry

	p, err := store.Get(id) //check
	assert.Equal(t, parcel, p)

	err = store.Delete(id) //deleting this db entry
	require.NoError(t, err)

	p, err = store.Get(id) //check
	require.Error(t, err)
	require.Empty(t, p)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel) //adding new parcel in db
	require.NoError(t, err)
	require.NotEmpty(t, id)

	newAddress := "new test address"
	err = store.SetAddress(id, newAddress) //setting new address
	require.NoError(t, err)

	p, err := store.Get(id) //check
	require.Equal(t, newAddress, p.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel) //adding new parcel in db
	require.NoError(t, err)
	require.NotEmpty(t, id)

	newStatus := ParcelStatusSent
	err = store.SetStatus(id, newStatus) //setting new status and checking for errors
	require.NoError(t, err)

	p, err := store.Get(id)
	require.Equal(t, newStatus, p.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotEmpty(t, id)

		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Equal(t, len(parcels), len(storedParcels))

	// check
	for _, parcel := range storedParcels {
		if parcel, exists := parcelMap[parcel.Number]; exists {
			require.Equal(t, parcel, parcelMap[parcel.Number])
		} else {
			panic("Посылка отсутствует в мапе")
		}
	}
}
