: Build the Project

go build -o task-cli task-cli.go

1 Add a task

./task-cli add "Buy groceries"

2 Update a task

./task-cli update 1 "Buy groceries and cook dinner"

3 Delete a task

./task-cli delete 1

.4 Mark as in-progress or done

./task-cli mark-in-progress 2
./task-cli mark-done 2


5 List tasks

./task-cli list           # all tasks
./task-cli list todo      # only pending
./task-cli list done      # only completed
./task-cli list in-progress  # in progress


6: Verify Task Storage

cat tasks.json
