package task

type TaskService struct {
	Repository TaskRepository
}

func NewService(repository TaskRepository) TaskService {
	return TaskService{
		Repository: repository,
	}
}
