import http from 'k6/http';
import {Rate} from 'k6/metrics';
import {generateTaskData} from './seeder.js';

const baseUrl = __ENV.HOSTNAME
const updateTaskFailRate = new Rate('failed form task update');

const urls = {
    tasksUrl: `${baseUrl}/tasks`,
};

export const createTask = (auth) => {
    const generatedTask = generateTaskData();
    const response = http.post(urls.tasksUrl, JSON.stringify(generatedTask), auth);
    return response.json().id
}

export const updateTask = (id, auth) => {
    const updatedTask = generateTaskData();

    const response = http.put(`${urls.tasksUrl}/${id}`, JSON.stringify(updatedTask), auth);
    updateTaskFailRate.add(response.status !== 200);

    return response.json().id
}

export const deleteTask = (id, auth) => {
     http.del(`${urls.tasksUrl}/${id}`, null, auth);
}
