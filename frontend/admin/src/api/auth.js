import request from '@/utils/request'

export function adminLogin(data) {
  return request({
    url: '/admin/login',
    method: 'post',
    data
  })
}

export function doctorLogin(data) {
  return request({
    url: '/doctor/login',
    method: 'post',
    data
  })
}

export function getStatistics() {
  return request({
    url: '/admin/statistics',
    method: 'get'
  })
}

export function getAppointmentStatistics(params) {
  return request({
    url: '/admin/statistics/appointments',
    method: 'get',
    params
  })
}
