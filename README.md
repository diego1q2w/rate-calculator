# Fare Estimator

### Intro
The flow is divided into 3 layers each one with their set of tests, those layers are `api` (the input), `app` 
 where all the calculations happen this one has more layers, and the `output` which as the name tells is the 
 the latest stage of the process, where the file is written.

- `api`: Is where the process starts with a property called `FileReader` this reads the provided file, validates its existence
and transforms each row into a go struct `Position` for further processing.

- `app`: Here is where the magic happens (and is divided into more stages), it accepts the struct provider by the `api` and processes 
them taking the advantage of goroutines, how those routines are used it depends on the stage (more about that below), and those stages are:
    
    - `segmenter`: This has the task to calculate a segment and the data associated with it, such as `distance`, `velocity` and `time`.
    - `filter` : This takes the output of the previous stages and filters it, in this case, based on speed.
    - `estimator`: This is the process that calculates the cost of that segment, based on mainly time and speed
    - `aggregator`: Uses the result of the estimator to aggregate the cost of every segment by ride ID

- `output`: It takes the result of the aggregator, witting into a file `output.txt`.

### Concurrency architecture

In this specific task, concurrency plays a very important role, as it is how the data is distributed across different `goroutines`
to accelerate its processing.

The `api` and `output` are single-threaded so nothing special about them, their job is to both read and write data from and to a single source.

The `app` however, is where the magic happens, each row is transformed into a struct called `Position`, then the `segmenter` uses pairs of them to
calculate a segment, that happens concurrently, the `segmenter` has `N` `goroutines` running then pushes those pairs of points to a channel, in each goroutine,
the segment is calculated (distance, time, etc), then is filtered, and the fare value of that segment is calculated.

In other words, each pair of Positions is routed using channels into a cluster of goroutines where within each goroutine, 
that pair is transformed into a segment (`segmenter`), then that segment is filtered (`filter`) 
and ultimately the fare value of that segment is calculated (`estimator`).

The aggregator runs in its cluster of `goroutines` which ideally should be less than the previous one, and the reason for that, 
is because each `goroutine` keeps its aggregation of the cost of the ride, so the result of the previous process (`estimator`), is then aggregated
in each goroutine independently from others, thanks to that we don't have to use a mutex and block the whole process each time a new increment is happening.

So each `aggregator` routine keeps its own and independent ride calculation, which means that a single ride can be aggregated in several 
`goroutines` at the same time. To merge those results we flush every now and then (with a channel), those workers into a `masterAggregator` 
which runs in his own `goroutine` and then process each flush one by one avoiding unwanted data overrides, to finally write that flush into a file

As the flush process might happen a few times in the process, the `output` file is updated each time.

### Tests

Each stage is unit-tested and the `app` has an end-to-end test, that would help how the app is bundled and how it may work
this file is located in `pkg/estimator/app/app_test.go`.

The `api` and `output` have their tests.

### How to run the project
If you have Docker composer installed you can use the scripts to run and test the project, for using a different path go to 
the `docker-compose.yml` and choose something else, by default uses `paths.csv` at the root of this project

For running:

    ./script/start
For testing:

    ./script/test
For linting:

    ./script/lint
    
Without Docker compose, at the root of this project run.

For running:

    go run main.go ./paths.csv  
For testing:

    go test ./...

After the script is done the result will be in a file called `output.txt` at the route of this project.
### Project structure

Each domain follow this convention
```
├── pkg/estimator
    ├─── app - It's where the business logic is placed.
       ├─── segmenter.go - It calculates each segment
       ├─── filter.go - Filters a given segment
       ├─── estimator-config.go - It has all the config needed by the estimator.go for calculating the fare segment
       ├─── estimator.go - Calculates the cost of that segment
       ├─── aggregator.go - Aggregates the cost of every segment by Ride ID
    ├─── domain - Are the domain entities.
    ├─── output - It's where we have the file writing logic.
    ├─── api - It's where we have the file reading logic.
├── script - Set of useful scripts, plug, and play.
    ├─── start
    ├─── lint
    ├─── test
└── ...
```
