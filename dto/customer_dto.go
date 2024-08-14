package dto

import "github.com/jacky-htg/billings/entity"

type CreateCustomerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func (u *CreateCustomerRequest) ToEntity() entity.Customer {
	return entity.Customer{
		Name:  u.Name,
		Phone: u.Phone,
		Email: u.Email,
	}
}

type CreateCustomerResponse struct {
	CustomerID int    `json:"customer_id"`
	Message    string `json:"message"`
}

func (resp *CreateCustomerResponse) FromEntity(e *entity.Customer) {
	resp.CustomerID = e.CustomerID
	resp.Message = "Customer created successfully"
}

type GetCustomerResponse struct {
	CustomerID int    `json:"customer_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
}

func (resp *GetCustomerResponse) FromEntity(e *entity.Customer) {
	resp.CustomerID = e.CustomerID
	resp.Name = e.Name
	resp.Email = e.Email
	resp.Phone = e.Phone
}
