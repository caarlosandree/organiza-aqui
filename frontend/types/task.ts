export interface TaskStatus {
  id: string
  user_id: string
  name: string
  color: string
  order_index: number
  is_default: boolean
  created_at: string
  updated_at: string
}

export interface Task {
  id: string
  user_id: string
  status_id: string
  title: string
  description: string
  priority: 'low' | 'medium' | 'high' | 'urgent'
  due_date?: string
  completed_at?: string
  lexorank: string
  financial_account_id?: string
  financial_amount?: number
  financial_category_id?: string
  created_at: string
  updated_at: string
}

export interface CreateTaskStatusRequest {
  name: string
  color: string
  order_index: number
  is_default: boolean
}

export interface UpdateTaskStatusRequest {
  name: string
  color: string
  order_index: number
  is_default: boolean
}

export interface CreateTaskRequest {
  status_id: string
  title: string
  description?: string
  priority: 'low' | 'medium' | 'high' | 'urgent'
  due_date?: string
}

export interface UpdateTaskRequest {
  status_id: string
  title: string
  description?: string
  priority: 'low' | 'medium' | 'high' | 'urgent'
  due_date?: string
}

export interface ReorderTaskRequest {
  task_id: string
  status_id: string
  after_id?: string
}

export interface TaskFilters {
  status_id?: string
  priority?: 'low' | 'medium' | 'high' | 'urgent'
  start_date?: string
  end_date?: string
  completed?: boolean
  limit?: number
  offset?: number
}
