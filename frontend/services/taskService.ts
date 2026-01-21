import api from '@/lib/axios'
import type {
  Task,
  TaskStatus,
  CreateTaskStatusRequest,
  UpdateTaskStatusRequest,
  CreateTaskRequest,
  UpdateTaskRequest,
  ReorderTaskRequest,
  TaskFilters,
} from '@/types/task'

export const getTaskStatuses = async (): Promise<TaskStatus[]> => {
  const response = await api.get<TaskStatus[]>('/task-statuses')
  return response.data
}

export const getTaskStatus = async (id: string): Promise<TaskStatus> => {
  const response = await api.get<TaskStatus>(`/task-statuses/${id}`)
  return response.data
}

export const createTaskStatus = async (
  data: CreateTaskStatusRequest
): Promise<TaskStatus> => {
  const response = await api.post<TaskStatus>('/task-statuses', data)
  return response.data
}

export const updateTaskStatus = async (
  id: string,
  data: UpdateTaskStatusRequest
): Promise<TaskStatus> => {
  const response = await api.put<TaskStatus>(`/task-statuses/${id}`, data)
  return response.data
}

export const deleteTaskStatus = async (id: string): Promise<void> => {
  await api.delete(`/task-statuses/${id}`)
}

export const reorderTaskStatuses = async (
  statusIds: string[]
): Promise<void> => {
  await api.post('/task-statuses/reorder', statusIds)
}

export const getTasks = async (filters?: TaskFilters): Promise<Task[]> => {
  const params = new URLSearchParams()
  if (filters?.status_id) params.append('status_id', filters.status_id)
  if (filters?.priority) params.append('priority', filters.priority)
  if (filters?.completed !== undefined)
    params.append('completed', String(filters.completed))
  if (filters?.limit) params.append('limit', String(filters.limit))
  if (filters?.offset) params.append('offset', String(filters.offset))

  const response = await api.get<Task[]>(`/tasks?${params.toString()}`)
  return response.data
}

export const getTask = async (id: string): Promise<Task> => {
  const response = await api.get<Task>(`/tasks/${id}`)
  return response.data
}

export const createTask = async (data: CreateTaskRequest): Promise<Task> => {
  const response = await api.post<Task>('/tasks', data)
  return response.data
}

export const updateTask = async (
  id: string,
  data: UpdateTaskRequest
): Promise<Task> => {
  const response = await api.put<Task>(`/tasks/${id}`, data)
  return response.data
}

export const deleteTask = async (id: string): Promise<void> => {
  await api.delete(`/tasks/${id}`)
}

export const reorderTask = async (
  data: ReorderTaskRequest
): Promise<Task> => {
  const response = await api.post<Task>('/tasks/reorder', data)
  return response.data
}

export const completeTask = async (id: string): Promise<Task> => {
  const response = await api.post<Task>(`/tasks/${id}/complete`)
  return response.data
}

export const uncompleteTask = async (id: string): Promise<Task> => {
  const response = await api.post<Task>(`/tasks/${id}/uncomplete`)
  return response.data
}
