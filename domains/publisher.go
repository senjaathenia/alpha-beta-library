package domains

//Publisher struct represents the Author model
type Publisher struct{
	ID			uint `gorm:"primaryKey" json:"id"`
	Name 		string `json: "name"`  // JSON tag harus di awali dengan json:
}

//PublisherRepository defines the methods that a repository
type PublisherRepository interface{
	Create(publisher *Publisher) error
	Update(publisher *Publisher) error
	Delete(id uint) error
	GetByID(id uint) (*Publisher, error)
	GetAll() ([]Publisher, error)
}

//PublisherUsecase defines the methods for the business logic
type PublisherUsecase interface{
	Create(publisher *Publisher) error
	Update(publisher *Publisher) error
	Delete(id uint) error
	GetByID(id uint) (*Publisher, error)
	GetAll() ([]Publisher, error)
}