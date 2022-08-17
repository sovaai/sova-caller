import axios from 'axios'

import { StartCampaignPayload, StopCampaignPayload } from '../@types/common'
import { API_URL } from '../constants/constants'

const api = axios.create({
  baseURL: API_URL,
})

const startCampaign = async (payload: StartCampaignPayload): Promise<string> => {
  try {
    const res = await api.post('/campaign/start', payload)
    return res.data
  } catch (error) {
    console.log(`Can't start campaign! Error: ${error}`)
    return ''
  }
}

const stopCampaign = async (payload: StopCampaignPayload): Promise<boolean> => {
  try {
    await api.post('/campaign/stop', payload)
    return true
  } catch (error) {
    console.log(`Can't stop campaign! Error: ${error}`)
    return false
  }
}

const getCampaignStatus = async (clientId: string) => {
  try {
    const res = await api.get(`/status/${clientId}`)
    return res.data
  } catch (error) {
    console.log(`Can't get campaign status! Error: ${error}`)
    return ''
  }
}

export { startCampaign, stopCampaign, getCampaignStatus }
