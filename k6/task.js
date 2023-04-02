import http from 'k6/http';
import {Rate} from 'k6/metrics';
import {generateTaskData} from './seeder.js';

const updateTaskFailRate = new Rate('failed_task_update');

const baseUrl = __ENV.HOSTNAME

export const createTask = (auth) => {
    const generatedTask = generateTaskData();
    const response = http.post(`${baseUrl}/tasks`, JSON.stringify(generatedTask), auth);
    return response.json().id
}

export const updateTask = (id, auth) => {
    const updatedTask = generateTaskData();

    const response = http.put(`${baseUrl}/tasks/${id}`, JSON.stringify(updatedTask), auth);
    updateTaskFailRate.add(response.status !== 200);

    return response.json().id
}

export const deleteTask = (id, auth) => {
     http.del(`${baseUrl}/tasks/${id}`, null, auth);
}
