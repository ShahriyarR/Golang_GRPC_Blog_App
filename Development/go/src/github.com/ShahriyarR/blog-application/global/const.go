package global

const (
	dburi = "mongodb+srv://standard:example@cluster0.ph8nx.mongodb.net/<dbname>?retryWrites=true&w=majority"
	dbname = "blog-application"
	performance = 100
)

var (
	jwtSecret = []byte("blogSecret")
)