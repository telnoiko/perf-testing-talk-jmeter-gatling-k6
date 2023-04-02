package tasks.model;

import org.apache.commons.lang3.RandomStringUtils;

public class Task {

    Integer id;
    String description;
    Boolean completed;

    public Task() {
    }

    public Task(Integer id, String description, Boolean completed) {
        this.id = id;
        this.description = description;
        this.completed = completed;
    }

    public static Task generateTask() {
        String description = RandomStringUtils.randomAlphabetic(10);
        Boolean completed = Math.random() < 0.5;
        Task t = new Task(null, description, completed);
        return t;
    }

    public Integer getId() {
        return id;
    }

    public void setId(Integer id) {
        this.id = id;
    }

    public String getDescription() {
        return description;
    }

    public Boolean getCompleted() {
        return completed;
    }
}
