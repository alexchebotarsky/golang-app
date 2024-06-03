package article

import "testing"

func TestPayload_Validate(t *testing.T) {
	type fields struct {
		Title       string
		Description string
		Body        string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "empty title",
			fields:  fields{Title: "", Description: "description", Body: "body"},
			wantErr: true,
		},
		{
			name:    "empty description",
			fields:  fields{Title: "title", Description: "", Body: "body"},
			wantErr: true,
		},
		{
			name:    "empty body",
			fields:  fields{Title: "title", Description: "description", Body: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Payload{
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Body:        tt.fields.Body,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Payload.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
