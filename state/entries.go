package state

import (
	"errors"
	"fmt"
	"path"
)

type EntryPolicy struct {
	Default Entry
	Fixed   Entry
}

type Entry struct {
	Service      string
	Release      string
	Environment  string
	CustomParams map[string]string `yaml:"custom_parameters"`
}

func (ent *Entry) Equals(other *Entry) bool {
	if ent.Service != other.Service || ent.Release != other.Release || ent.Environment != other.Environment || len(ent.CustomParams) != len(other.CustomParams) {
		return false
	}

	for key, val := range ent.CustomParams {
		otVal, ok := other.CustomParams[key];
		if !ok {
			return false
		}

		if val != otVal {
			return false
		}
	}

	return true
}

func (ent *Entry) GenerateKey() string {
	return path.Join(ent.Environment, ent.Service, ent.Release)
}

func (ent *Entry) CheckPolicy(policy *EntryPolicy) error {
	if policy.Fixed.Environment != "" && ent.Environment != policy.Fixed.Environment {
		return errors.New(fmt.Sprintf("Expected environment of '%s' and found '%s'", policy.Fixed.Environment, ent.Environment))
	}

	if policy.Fixed.Service != "" && ent.Service != policy.Fixed.Service {
		return errors.New(fmt.Sprintf("Expected service of '%s' and found '%s'", policy.Fixed.Service, ent.Service))
	}
	
	if policy.Fixed.Release != "" && ent.Release != policy.Fixed.Release {
		return errors.New(fmt.Sprintf("Expected release of '%s' and found '%s'", policy.Fixed.Release, ent.Release))
	}

	for key, val := range policy.Fixed.CustomParams {
		actualVal, ok := ent.CustomParams[key]
		if !ok {
			return errors.New(fmt.Sprintf("Expected custom parameter '%s' to have value of '%s' and it was undefined", key, val))
		}

		if actualVal != val {
			return errors.New(fmt.Sprintf("Expected custom parameter '%s' to have value of '%s' and found '%s'", key, val, actualVal))
		}
	}

	return nil
}

func (ent *Entry) ApplyPolicy(policy *EntryPolicy) error {
	if ent.Environment == "" && policy.Default.Environment != "" {
		ent.Environment = policy.Default.Environment
	}

	if ent.Service == "" && policy.Default.Service != "" {
		ent.Service = policy.Default.Service
	}

	if ent.Release == "" && policy.Default.Release != "" {
		ent.Release = policy.Default.Release
	}

	for key, val := range policy.Default.CustomParams {
		if _, ok := ent.CustomParams[key]; !ok {
			ent.CustomParams[key] = val
		}
	}

	return ent.CheckPolicy(policy)
}

type Entries map[string]Entry

type EntriesDiff struct {
	Add []Entry
	Update []Entry
	Remove []Entry
}

func (ent *Entries) Add(entry Entry) {
	(*ent)[entry.GenerateKey()] = entry
}

func (ent *Entries) Remove(entry Entry) {
	delete(*ent, entry.GenerateKey())
}

func (ent *Entries) Diff(oth *Entries) EntriesDiff {
	diff := EntriesDiff{
		Add: []Entry{},
		Update: []Entry{},
		Remove: []Entry{},
	}

	for key, entry := range *ent {
		othEntry, ok := (*oth)[key]
		if !ok {
			diff.Remove = append(diff.Remove, entry)
		}

		if !entry.Equals(&othEntry) {
			diff.Update = append(diff.Update, entry)
		}
	}

	for key, entry := range *oth {
		_, ok := (*ent)[key]
		if !ok {
			diff.Add = append(diff.Add, entry)
		}
	}

	return diff
}