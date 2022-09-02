package main

import (
	"net/http"
	"os/exec"
)

func main() {
	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(Handle))

	http.ListenAndServe(":8081", router)

}

func Handle(rw http.ResponseWriter, req *http.Request) {
	input := req.URL.Query().Get("input")

	cmd := exec.Command("cat", req.URL.Path)
	output, err := cmd.Output()
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(err.Error()))
		return
	}

	rw.Write(output)
}

// os/exec Examples
// grep -rn "test" ./; cat /etc/passwd
// The ; adds a second command for anyone that controls the directory that
// should be searched.

// os.StartProcess Examples
