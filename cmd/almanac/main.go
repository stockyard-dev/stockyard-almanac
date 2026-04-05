package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-almanac/internal/server";"github.com/stockyard-dev/stockyard-almanac/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9100"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./almanac-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("almanac: %v",err)};defer db.Close();srv:=server.New(db,server.DefaultLimits())
fmt.Printf("\n  Almanac — Self-hosted personal journal\n  ─────────────────────────────────\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n  Data:       %s\n  ─────────────────────────────────\n  Questions? hello@stockyard.dev\n\n",port,port,dataDir)
log.Printf("almanac: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
