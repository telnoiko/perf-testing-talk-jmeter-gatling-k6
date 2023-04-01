package tasks.model;

public class Task {

    String description;
    Boolean completed;

    public Task() {
    }

    public Task(String description, Boolean completed) {
        this.description = description;
        this.completed = completed;
    }
}
