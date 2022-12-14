# Purpose

Run tests performed with curl

## Description

The runt tests for accessing the API with curl. The run test normally
has first option: number of iterations which helps to do some performance
testing runs. The test scope is randomized ie. the size of the file
uploaded is random.

## Test data

The temporary test data use the folder: ./testdata to store the files used
during run test.

## Results

The .lst files with results of each test step are stored in the folder: ./testresults.
Usually they contain the return status of each curl call with the printout
of the result, usually formatted as json.