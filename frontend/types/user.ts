// User.ts - 用户相关类型定义

/**
 * 用户实体接口
 */
export interface User {
  id: number;
  username: string;
  email: string;
  nickname: string;
  avatar_url: string;
  bio: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

/**
 * 用户注册请求参数
 */
export interface RegisterRequest {
  username: string;
  password: string;
  email: string;
  nickname: string;
}

/**
 * 用户登录请求参数
 */
export interface LoginRequest {
  username: string;
  password: string;
}

/**
 * 用户资料更新请求参数
 */
export interface UpdateProfileRequest {
  nickname: string;
  avatar_url: string;
  bio: string;
}