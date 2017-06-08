package datatypes
import (
	"net/http"
)

type GraphqlServer struct {
	Running bool
	Instance *http.Server
}
