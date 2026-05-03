import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/user/login',
    method: 'post',
    data
  })
}

export function register(data) {
  return request({
    url: '/user/register',
    method: 'post',
    data
  })
}

export function getUserInfo() {
  return request({
    url: '/user/info',
    method: 'get'
  })
}

export function updateUserInfo(data) {
  return request({
    url: '/user/info',
    method: 'put',
    data
  })
}

export function updatePassword(data) {
  return request({
    url: '/user/password',
    method: 'put',
    data
  })
}

export function getUserAppointments(params) {
  return request({
    url: '/user/appointments',
    method: 'get',
    params
  })
}

export function createAppointment(data) {
  return request({
    url: '/user/appointments',
    method: 'post',
    data
  })
}

export function cancelAppointment(id) {
  return request({
    url: `/user/appointments/${id}`,
    method: 'delete'
  })
}

export function getUserConsultations(params) {
  return request({
    url: '/user/consultations',
    method: 'get',
    params
  })
}

export function commentConsultation(id, data) {
  return request({
    url: `/user/consultations/${id}/comment`,
    method: 'post',
    data
  })
}

export function getFavorites(params) {
  return request({
    url: '/user/favorites',
    method: 'get',
    params
  })
}

export function addFavorite(data) {
  return request({
    url: '/user/favorites',
    method: 'post',
    data
  })
}

export function removeFavorite(id) {
  return request({
    url: `/user/favorites/${id}`,
    method: 'delete'
  })
}
