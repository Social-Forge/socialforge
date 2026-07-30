package entity

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NullString struct {
	sql.NullString
}
type NullTime struct {
	sql.NullTime
}

type NullInt16 struct {
	sql.NullInt16
}
type NullInt32 struct {
	sql.NullInt32
}
type NullInt64 struct {
	sql.NullInt64
}
type NullFloat64 struct {
	sql.NullFloat64
}
type NullBool struct {
	sql.NullBool
}
type NullByte struct {
	sql.NullByte
}
type StringArray struct {
	Elements []string
	Valid    bool
}
type Int32Array struct {
	Elements []int32
	Valid    bool
}
type Int64Array struct {
	Elements []int64
	Valid    bool
}
type Float64Array struct {
	Elements []float64
	Valid    bool
}
type BoolArray struct {
	Elements []bool
	Valid    bool
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}
func (ns *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		return nil
	}
	ns.Valid = true
	return json.Unmarshal(data, &ns.String)
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time)
}
func (nt *NullTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nt.Valid = false
		return nil
	}
	nt.Valid = true
	return json.Unmarshal(data, &nt.Time)
}

func (ni NullInt16) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int16)
}
func (ni *NullInt16) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ni.Valid = false
		return nil
	}
	ni.Valid = true
	return json.Unmarshal(data, &ni.Int16)
}

func (ni NullInt32) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int32)
}
func (ni *NullInt32) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ni.Valid = false
		return nil
	}
	ni.Valid = true
	return json.Unmarshal(data, &ni.Int32)
}

func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}
func (ni *NullInt64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ni.Valid = false
		return nil
	}
	ni.Valid = true
	return json.Unmarshal(data, &ni.Int64)
}

func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Float64)
}
func (nf *NullFloat64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nf.Valid = false
		return nil
	}
	nf.Valid = true
	return json.Unmarshal(data, &nf.Float64)
}

func (nb NullBool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nb.Bool)
}
func (nb *NullBool) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nb.Valid = false
		return nil
	}
	nb.Valid = true
	return json.Unmarshal(data, &nb.Bool)
}

func (nb NullByte) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nb.Byte)
}
func (nb *NullByte) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nb.Valid = false
		return nil
	}
	nb.Valid = true
	return json.Unmarshal(data, &nb.Byte)
}

func (sa StringArray) MarshalJSON() ([]byte, error) {
	if !sa.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(sa.Elements)
}
func (sa *StringArray) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		sa.Valid = false
		return nil
	}
	sa.Valid = true
	return json.Unmarshal(data, &sa.Elements)
}

func (ia Int32Array) MarshalJSON() ([]byte, error) {
	if !ia.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ia.Elements)
}
func (ia *Int32Array) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ia.Valid = false
		return nil
	}
	ia.Valid = true
	return json.Unmarshal(data, &ia.Elements)
}

func (ia Int64Array) MarshalJSON() ([]byte, error) {
	if !ia.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ia.Elements)
}
func (ia *Int64Array) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ia.Valid = false
		return nil
	}
	ia.Valid = true
	return json.Unmarshal(data, &ia.Elements)
}

func (fa Float64Array) MarshalJSON() ([]byte, error) {
	if !fa.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(fa.Elements)
}
func (fa *Float64Array) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		fa.Valid = false
		return nil
	}
	fa.Valid = true
	return json.Unmarshal(data, &fa.Elements)
}

func NewNullString(s string) NullString {
	if s == "" {
		return NullString{NullString: sql.NullString{Valid: false}}
	}
	return NullString{NullString: sql.NullString{String: s, Valid: true}}
}
func NewNullStringPtr(s *string) NullString {
	if s == nil {
		return NullString{NullString: sql.NullString{Valid: false}}
	}
	return NullString{NullString: sql.NullString{String: *s, Valid: true}}
}

func NewNullTime(t time.Time) NullTime {
	if t.IsZero() {
		return NullTime{NullTime: sql.NullTime{Valid: false}}
	}
	return NullTime{NullTime: sql.NullTime{Time: t, Valid: true}}
}
func NewNullTimePtr(t *time.Time) NullTime {
	if t == nil || t.IsZero() {
		return NullTime{NullTime: sql.NullTime{Valid: false}}
	}
	return NullTime{NullTime: sql.NullTime{Time: *t, Valid: true}}
}
func NewNullTimeFromNow(years, months, days int) NullTime {
	futureTime := time.Now().AddDate(years, months, days)
	return NullTime{NullTime: sql.NullTime{Time: futureTime, Valid: true}}
}
func NewNullTimeFromDate(year, month, day int) NullTime {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return NullTime{NullTime: sql.NullTime{Time: t, Valid: true}}
}

func NewNullInt16(i int16) NullInt16 {
	return NullInt16{NullInt16: sql.NullInt16{Int16: i, Valid: true}}
}
func NewNullInt16Ptr(i *int16) NullInt16 {
	if i == nil {
		return NullInt16{NullInt16: sql.NullInt16{Valid: false}}
	}
	return NullInt16{NullInt16: sql.NullInt16{Int16: *i, Valid: true}}
}

func NewNullInt32(i int32) NullInt32 {
	return NullInt32{NullInt32: sql.NullInt32{Int32: i, Valid: true}}
}
func NewNullInt32Ptr(i *int32) NullInt32 {
	if i == nil {
		return NullInt32{NullInt32: sql.NullInt32{Valid: false}}
	}
	return NullInt32{NullInt32: sql.NullInt32{Int32: *i, Valid: true}}
}

func NewNullInt64(i int64) NullInt64 {
	return NullInt64{NullInt64: sql.NullInt64{Int64: i, Valid: true}}
}
func NewNullInt64Ptr(i *int64) NullInt64 {
	if i == nil {
		return NullInt64{NullInt64: sql.NullInt64{Valid: false}}
	}
	return NullInt64{NullInt64: sql.NullInt64{Int64: *i, Valid: true}}
}

func NewNullFloat64(f float64) NullFloat64 {
	return NullFloat64{NullFloat64: sql.NullFloat64{Float64: f, Valid: true}}
}
func NewNullFloat64Ptr(f *float64) NullFloat64 {
	if f == nil {
		return NullFloat64{NullFloat64: sql.NullFloat64{Valid: false}}
	}
	return NullFloat64{NullFloat64: sql.NullFloat64{Float64: *f, Valid: true}}
}

func NewNullBool(b bool) NullBool {
	return NullBool{NullBool: sql.NullBool{Bool: b, Valid: true}}
}
func NewNullBoolPtr(b *bool) NullBool {
	if b == nil {
		return NullBool{NullBool: sql.NullBool{Valid: false}}
	}
	return NullBool{NullBool: sql.NullBool{Bool: *b, Valid: true}}
}

func NewNullByte(b byte) NullByte {
	return NullByte{NullByte: sql.NullByte{Byte: b, Valid: true}}
}
func NewNullBytePtr(b *byte) NullByte {
	if b == nil {
		return NullByte{NullByte: sql.NullByte{Valid: false}}
	}
	return NullByte{NullByte: sql.NullByte{Byte: *b, Valid: true}}
}

func NewStringArray(elements []string) StringArray {
	if elements == nil {
		return StringArray{Valid: false}
	}
	return StringArray{Elements: elements, Valid: true}
}
func NewStringArrayPtr(elements *[]string) StringArray {
	if elements == nil {
		return StringArray{Valid: false}
	}
	return StringArray{Elements: *elements, Valid: true}
}

func NewInt32Array(elements []int32) Int32Array {
	if elements == nil {
		return Int32Array{Valid: false}
	}
	return Int32Array{Elements: elements, Valid: true}
}
func NewInt32ArrayPtr(elements *[]int32) Int32Array {
	if elements == nil {
		return Int32Array{Valid: false}
	}
	return Int32Array{Elements: *elements, Valid: true}
}

func NewInt64Array(elements []int64) Int64Array {
	if elements == nil {
		return Int64Array{Valid: false}
	}
	return Int64Array{Elements: elements, Valid: true}
}
func NewInt64ArrayPtr(elements *[]int64) Int64Array {
	if elements == nil {
		return Int64Array{Valid: false}
	}
	return Int64Array{Elements: *elements, Valid: true}
}

func NewFloat64Array(elements []float64) Float64Array {
	if elements == nil {
		return Float64Array{Valid: false}
	}
	return Float64Array{Elements: elements, Valid: true}
}
func NewFloat64ArrayPtr(elements *[]float64) Float64Array {
	if elements == nil {
		return Float64Array{Valid: false}
	}
	return Float64Array{Elements: *elements, Valid: true}
}

func (nt NullTime) GetTime() *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}
func (ns NullString) GetString() *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}
func (ni NullInt16) GetInt16() *int16 {
	if ni.Valid {
		return &ni.Int16
	}
	return nil
}
func (ni NullInt32) GetInt32() *int32 {
	if ni.Valid {
		return &ni.Int32
	}
	return nil
}
func (ni NullInt64) GetInt64() *int64 {
	if ni.Valid {
		return &ni.Int64
	}
	return nil
}
func (nf NullFloat64) GetFloat64() *float64 {
	if nf.Valid {
		return &nf.Float64
	}
	return nil
}
func (nb NullBool) GetBool() *bool {
	if nb.Valid {
		return &nb.Bool
	}
	return nil
}
func (nb NullByte) GetByte() *byte {
	if nb.Valid {
		return &nb.Byte
	}
	return nil
}
func (sa StringArray) GetElements() *[]string {
	if sa.Valid {
		return &sa.Elements
	}
	return nil
}
func (ia Int32Array) GetElements() *[]int32 {
	if ia.Valid {
		return &ia.Elements
	}
	return nil
}
func (ia Int64Array) GetElements() *[]int64 {
	if ia.Valid {
		return &ia.Elements
	}
	return nil
}
func (fa Float64Array) GetElements() *[]float64 {
	if fa.Valid {
		return &fa.Elements
	}
	return nil
}

func NullInt16FromDB(i sql.NullInt16) NullInt16 {
	return NullInt16{NullInt16: i}
}
func NullInt32FromDB(i sql.NullInt32) NullInt32 {
	return NullInt32{NullInt32: i}
}
func NullInt64FromDB(i sql.NullInt64) NullInt64 {
	return NullInt64{NullInt64: i}
}
func NullFloat64FromDB(f sql.NullFloat64) NullFloat64 {
	return NullFloat64{NullFloat64: f}
}
func NullBoolFromDB(b sql.NullBool) NullBool {
	return NullBool{NullBool: b}
}
func NullByteFromDB(b sql.NullByte) NullByte {
	return NullByte{NullByte: b}
}
