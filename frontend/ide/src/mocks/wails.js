console.log('[Mock] Wails bindings loaded');

export function SelectWorkspace() {
    return Promise.resolve("mock-workspace-id");
}

export function ListTasks() {
    console.log("[Mock] ListTasks called");
    // Return mock tasks
    const tasks = JSON.parse(window.localStorage.getItem('mock_tasks') || '[]');
    return Promise.resolve(tasks);
}

export function GetWorkspace(id) {
    return Promise.resolve({
        version: "1.0",
        projectRoot: "/mock/root",
        displayName: "Mock Project"
    });
}

export function CreateTask(title, poolId) {
    console.log("[Mock] CreateTask called", title, poolId);
    const tasks = JSON.parse(window.localStorage.getItem('mock_tasks') || '[]');
    const newTask = {
        id: "task-" + Date.now(),
        title: title,
        status: "PENDING",
        poolId: poolId,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString()
    };
    tasks.push(newTask);
    window.localStorage.setItem('mock_tasks', JSON.stringify(tasks));
    return Promise.resolve(newTask);
}

export function RunTask(taskId) {
    console.log("[Mock] RunTask called", taskId);
    return Promise.resolve();
}
