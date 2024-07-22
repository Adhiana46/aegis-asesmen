package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	Config "github.com/Adhiana46/aegis-asesmen/config"
	Migrations "github.com/Adhiana46/aegis-asesmen/database/migrations"
	PkgDataSource "github.com/Adhiana46/aegis-asesmen/pkg/data_sources"
	PkgKafkaPublisher "github.com/Adhiana46/aegis-asesmen/pkg/kafka/publisher"
	PkgHttpServer "github.com/Adhiana46/aegis-asesmen/pkg/server/http"

	UserDataSourcePostgres "github.com/Adhiana46/aegis-asesmen/internal/user/data/source/postgres"
	UserHttpHandler "github.com/Adhiana46/aegis-asesmen/internal/user/delivery/http_handler"
	UserRepository "github.com/Adhiana46/aegis-asesmen/internal/user/domain/repository"
	UserUsecase "github.com/Adhiana46/aegis-asesmen/internal/user/domain/usecase"

	OrganizationDataSourcePostgres "github.com/Adhiana46/aegis-asesmen/internal/organization/data/source/postgres"
	OrganizationHttpHandler "github.com/Adhiana46/aegis-asesmen/internal/organization/delivery/http_handler"
	OrganizationRepository "github.com/Adhiana46/aegis-asesmen/internal/organization/domain/repository"
	OrganizationUsecase "github.com/Adhiana46/aegis-asesmen/internal/organization/domain/usecase"

	"github.com/IBM/sarama"
	"github.com/golang-migrate/migrate"
	"github.com/pkg/errors"
	"github.com/xdg-go/scram"
)

func main() {
	cfg, err := Config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// Init DB
	dbConn, err := initPostgres(cfg)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := dbConn.Close(); err != nil {
			slog.Error("error during cleanup persistent", slog.String("error", err.Error()))
		}
	}()

	// Init Publisher
	// publisher, err := initPublisher(cfg)
	// if err != nil {
	// 	panic(err)
	// }
	// defer func() {
	// 	if err := publisher.Close(); err != nil {
	// 		slog.Error("error during cleanup publisher", slog.String("error", err.Error()))
	// 	}
	// }()

	// Data Source
	var (
		userPersistent         = UserDataSourcePostgres.NewUserPersistentPostgres(dbConn)
		organizationPersistent = OrganizationDataSourcePostgres.NewOrganizationPersistentPostgres(dbConn)
	)

	// Repositories
	var (
		userRepository         = UserRepository.NewUserRepository(userPersistent)
		organizationRepository = OrganizationRepository.NewOrganizationRepository(organizationPersistent)
	)

	// Init Usecase
	var (
		userSigninUsecase          = UserUsecase.NewUserSigninUsecase(cfg, userRepository)
		getListOrganizationUsecase = OrganizationUsecase.NewGetListOrganizationUsecase(cfg, organizationRepository)
		getOrganizationByIdUsecase = OrganizationUsecase.NewGetOrganizationByIdUsecase(cfg, organizationRepository)
		createOrganizationUsecase  = OrganizationUsecase.NewCreateOrganizationUsecase(cfg, organizationRepository)
		updateOrganizationUsecase  = OrganizationUsecase.NewUpdateOrganizationUsecase(cfg, organizationRepository)
		deleteOrganizationUsecase  = OrganizationUsecase.NewDeleteOrganizationUsecase(cfg, organizationRepository)
	)

	// Init Handlers
	var (
		userHttpHandler         = UserHttpHandler.NewAuthHandler(cfg, userSigninUsecase)
		organizationHttpHandler = OrganizationHttpHandler.NewOrganizationHandler(
			cfg,
			getListOrganizationUsecase,
			getOrganizationByIdUsecase,
			createOrganizationUsecase,
			updateOrganizationUsecase,
			deleteOrganizationUsecase,
		)
	)

	// Http Server
	httpServer := PkgHttpServer.New(
		PkgHttpServer.Address(cfg.HttpServer.Host, cfg.HttpServer.Port),
		PkgHttpServer.WithNameAndVersion(cfg.App.Name, cfg.App.Version),
	)

	// Setup Http Server
	setupMiddlewares(cfg, httpServer)
	setupValidator(cfg, httpServer)
	setupErrorHandler(cfg, httpServer)
	setupHttpRoutes(
		cfg,
		httpServer,
		// Handlers
		userHttpHandler,
		organizationHttpHandler,
	)

	// Start HTTP Server
	httpServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		slog.Info("Interrupted: " + s.String())
	case err = <-httpServer.Notify():
		slog.Error(fmt.Errorf("httpServer.Notify: %w", err).Error())
	}

	// Shutdown
	if err := httpServer.Shutdown(); err != nil {
		slog.Error(fmt.Errorf("httpServer.Shutdown: %w", err).Error())
	}
}

func initPostgres(cfg *Config.Config) (*PkgDataSource.PostgresDB, error) {
	dbConn := PkgDataSource.NewPostgresDb(
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Pass,
		cfg.Postgres.DbName,
	)

	// RUN Migrations
	dbMigrator := Migrations.NewPostgresMigrator(dbConn, "./database/migrations/postgres")
	if dbMigrator != nil {
		err := dbMigrator.Up()
		if err != nil && err.Error() != migrate.ErrNoChange.Error() {
			return nil, err
		}
	}

	return dbConn, nil
}

func initPublisher(cfg *Config.Config) (*PkgKafkaPublisher.Publisher, error) {
	config := sarama.NewConfig()
	config.ClientID = cfg.Kafka.ClientId
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	// SASL_SSL Configuration
	config.Net.SASL.Enable = true
	config.Net.SASL.Handshake = true
	config.Net.SASL.User = cfg.Kafka.Username
	config.Net.SASL.Password = cfg.Kafka.Password
	config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
	if cfg.Kafka.SASLMechanism == "SCRAM-SHA-512" {
		config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
	}
	config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
		var hashGenerator scram.HashGeneratorFcn = SHA256
		if cfg.Kafka.SASLMechanism == "SCRAM-SHA-512" {
			hashGenerator = SHA512
		}

		return &XDGSCRAMClient{HashGeneratorFcn: hashGenerator}
	}
	config.Net.TLS.Enable = true

	publisher, err := PkgKafkaPublisher.New(
		PkgKafkaPublisher.WithConfig(config),
		PkgKafkaPublisher.WithBrokers(cfg.Kafka.Brokers),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start kafka publisher")
	}

	return publisher, nil
}
