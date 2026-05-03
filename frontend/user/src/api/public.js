import request from '@/utils/request'

export function getCarousels() {
  return request({
    url: '/carousels',
    method: 'get'
  })
}

export function getAnnouncements(params) {
  return request({
    url: '/announcements',
    method: 'get',
    params
  })
}

export function getAnnouncementDetail(id) {
  return request({
    url: `/announcements/${id}`,
    method: 'get'
  })
}

export function getMedicines(params) {
  return request({
    url: '/medicines',
    method: 'get',
    params
  })
}

export function getMedicineDetail(id) {
  return request({
    url: `/medicines/${id}`,
    method: 'get'
  })
}

export function getMedicalProjects(params) {
  return request({
    url: '/medical-projects',
    method: 'get',
    params
  })
}

export function getMedicalProjectDetail(id) {
  return request({
    url: `/medical-projects/${id}`,
    method: 'get'
  })
}

export function getHospitals() {
  return request({
    url: '/hospitals',
    method: 'get'
  })
}

export function getDepartments(params) {
  return request({
    url: '/departments',
    method: 'get',
    params
  })
}

export function getDoctors(params) {
  return request({
    url: '/doctors',
    method: 'get',
    params
  })
}

export function getDoctorDetail(id) {
  return request({
    url: `/doctors/${id}`,
    method: 'get'
  })
}

export function getDoctorSchedules(id, params) {
  return request({
    url: `/doctors/${id}/schedules`,
    method: 'get',
    params
  })
}
