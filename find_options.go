package common

type FindOption func(*FindOptions)

type FindOptions struct {
	fieldEquals    map[string]interface{}
	fieldNotEquals map[string]interface{}
}

func NewFindOptions() *FindOptions {
	return &FindOptions{
		fieldEquals:    make(map[string]interface{}),
		fieldNotEquals: make(map[string]interface{}),
	}
}

func WithFieldEquals(field string, value interface{}) FindOption {
	return func(o *FindOptions) {
		o.fieldEquals[field] = value
	}
}

func WithFieldNotEquals(field string, value interface{}) FindOption {
	return func(o *FindOptions) {
		o.fieldNotEquals[field] = value
	}
}

func (o *FindOptions) FieldEquals() map[string]interface{} {
	return o.fieldEquals
}

func (o *FindOptions) FieldNotEquals() map[string]interface{} {
	return o.fieldNotEquals
}
