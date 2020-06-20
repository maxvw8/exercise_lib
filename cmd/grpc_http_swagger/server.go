package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/maxvw8/exercise_lib/exrs"
	"github.com/maxvw8/exercise_lib/exrs/storage/mongodb"
	pbexrs "github.com/maxvw8/exercise_lib/pbexrs/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func newServer() *exrs.API {
	//create repository connection
	repo, err := mongodb.New("myDB")
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()
	// Create BlogService type
	srv, err := exrs.Server(repo)
	if err != nil {
		log.Fatal(err)
	}
	return srv
}
func getKeyPair() tls.Certificate {
	cerPath, _ := filepath.Abs("certs/server.pem")
	keyPath, _ := filepath.Abs("certs/server.key")
	cer, err := tls.LoadX509KeyPair(cerPath, keyPath)
	if err != nil {
		log.Fatal(err)
	}
	return cer
}

func getCertPool(cer tls.Certificate) *x509.CertPool {
	demoCertPool := x509.NewCertPool()
	cert, _ := ioutil.ReadFile("certs/server.pem")
	ok := demoCertPool.AppendCertsFromPEM(cert)
	if !ok {
		log.Fatal("bad certs")
	}
	return demoCertPool
}
func main() {
	srvAddress, port := "localhost", "10000"
	srvAddress = fmt.Sprintf("%s:%s", srvAddress, port)
	keypair := getKeyPair()
	certPool := getCertPool(keypair)
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any

	// Set options, here we can configure things like TLS support
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(logger),
		)),
	}
	// Create new gRPC server with (blank) options
	grpcServer := grpc.NewServer(opts...)

	pbexrs.RegisterExerciseServiceServer(grpcServer, newServer())
	ctx := context.Background()
	dcreds := credentials.NewTLS(&tls.Config{
		ServerName: srvAddress,
		RootCAs:    certPool,
	})
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}

	mux := http.NewServeMux()
	gwmux := runtime.NewServeMux()

	err := pbexrs.RegisterExerciseServiceHandlerFromEndpoint(ctx, gwmux, srvAddress, dopts)
	if err != nil {
		fmt.Printf("serve: %v\n", err)
		return
	}
	mux.Handle("/", gwmux)

	//serve swagger
	conn, err := net.Listen("tcp", srvAddress)
	if err != nil {
		panic(err)
	}
	srv := &http.Server{
		Addr:    srvAddress,
		Handler: grpcHandlerFunc(grpcServer, mux),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{keypair},
			NextProtos:   []string{"h2"},
		},
	}
	fmt.Printf("grpc on port: %s\n", port)
	err = srv.Serve(conn)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
func grpcHandlerFunc(grpcServer *grpc.Server, httpHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			httpHandler.ServeHTTP(w, r)
		}
	})
}
