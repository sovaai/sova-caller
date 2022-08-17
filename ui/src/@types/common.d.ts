import { ChangeEvent } from 'react'

export type InputChangeEventHandlerType = (e: ChangeEvent<HTMLInputElement>) => void
export type TextareaChangeEventHandlerType = (e: ChangeEvent<HTMLTextAreaElement>) => void

export interface StartCampaignPayload {
  uuid: string
  capi: string
  contacts: string[]
  sip: {
    sip_registrar: string
    sip_id: string
    sip_password: string
  }
  asr: {
    url: string
    token: string
  }
  tts: {
    url: string
    token: string
  }
  id: string
}

export interface StopCampaignPayload {
  id: string
}
