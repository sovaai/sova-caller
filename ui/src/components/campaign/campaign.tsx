import { Button, Flex, Text } from '@chakra-ui/react'
import React, { memo, useCallback, useEffect, useRef, useState } from 'react'

import { StartCampaignPayload } from '../../@types/common'
import { getCampaignStatus, startCampaign, stopCampaign } from '../../api/api'
import { parseContacts } from '../../helpers/parseContacts'
import { useLocalStorage } from '../../hooks/useLocalStorage'
import { colors } from '../../styles/theme'
import { Contacts } from '../contacts/contacts'
import { Settings } from '../settings/settings'

const Campaign = memo(() => {
  const [capiUrl, setcapiUrl] = useLocalStorage('capiUrl', '')
  const [uuid, setUuid] = useLocalStorage('uuid', '')

  const [asrUrl, setAsrUrl] = useLocalStorage('asrUrl', '')
  const [asrToken, setAsrToken] = useLocalStorage('asrToken', '')

  const [ttsUrl, setTtsUrl] = useLocalStorage('ttsUrl', '')
  const [ttsToken, setTtsToken] = useLocalStorage('ttsToken', '')

  const [sipRegistrar, setSipRegistrar] = useLocalStorage('sipRegistrar', '')
  const [sipId, setSipId] = useLocalStorage('sipId', '')
  const [sipPassword, setSipPassword] = useLocalStorage('sipPassword', '')

  const [campaignStatus, setCampaignStatus] = useLocalStorage('campaignStatus', '')

  const [rawContacts, setRawContacts] = useLocalStorage('rawContacts', '')

  const [clientId, setClientId] = useLocalStorage('clientId', '')

  const [isLoading, setIsLoading] = useState(false)

  const isInitial = useRef(true)

  const updateStatus = useCallback(() => {
    const getStatus = async () => {
      const res = await getCampaignStatus(clientId)
      setCampaignStatus(res)

      if (res === 'running') {
        setTimeout(getStatus, 3000)
      }
    }

    setTimeout(getStatus, 2000)
  }, [clientId, setCampaignStatus])

  useEffect(() => {
    if (campaignStatus === 'running') {
      updateStatus()
    }
    setIsLoading(false)
  }, [campaignStatus, updateStatus])

  useEffect(() => {
    if (isInitial.current) {
      isInitial.current = false
      return
    }
    setCampaignStatus('')
  }, [
    capiUrl,
    uuid,
    asrUrl,
    asrToken,
    ttsUrl,
    ttsToken,
    sipRegistrar,
    sipId,
    sipPassword,
    rawContacts,
    setCampaignStatus,
  ])

  const buttonText = campaignStatus === 'running' ? 'Остановить обзвон' : 'Запустить обзвон'
  const disabled = campaignStatus === 'running'

  const assistantSettings = {
    heading: 'Подключение Ассистента',
    fields: [
      { value: capiUrl, setter: setcapiUrl, placeholder: 'URL' },
      { value: uuid, setter: setUuid, placeholder: 'UUID канала' },
    ],
    disabled,
  }

  const asrSettings = {
    heading: 'Распознавание речи ASR',
    fields: [
      { value: asrUrl, setter: setAsrUrl, placeholder: 'URL' },
      { value: asrToken, setter: setAsrToken, placeholder: 'Token' },
    ],
    disabled,
  }

  const ttsSettings = {
    heading: 'Синтез речи TTS',
    fields: [
      { value: ttsUrl, setter: setTtsUrl, placeholder: 'URL' },
      { value: ttsToken, setter: setTtsToken, placeholder: 'Token' },
    ],
    disabled,
  }

  const sipSettings = {
    heading: 'Телефония',
    fields: [
      { value: sipRegistrar, setter: setSipRegistrar, placeholder: 'Proxy-сервер' },
      { value: sipId, setter: setSipId, placeholder: 'SIP login' },
      { value: sipPassword, setter: setSipPassword, placeholder: 'SIP password' },
    ],
    disabled,
  }

  const isIncorrectSettings = [
    capiUrl,
    uuid,
    asrUrl,
    asrToken,
    ttsUrl,
    ttsToken,
    sipRegistrar,
    sipId,
    sipPassword,
    rawContacts,
  ].some((s) => s === '')

  const changeStatus = async () => {
    if (campaignStatus === 'running') {
      setIsLoading(true)
      const res = await stopCampaign({
        id: clientId,
      })

      if (!res) {
        setCampaignStatus('failed')
        setIsLoading(false)
      }

      return
    }

    const contacts: string[] = parseContacts(rawContacts)

    if (contacts.length === 0) {
      setCampaignStatus('failed')
      return
    }

    const payload: StartCampaignPayload = {
      uuid,
      capi: capiUrl,
      contacts,
      sip: {
        sip_registrar: sipRegistrar,
        sip_id: sipId,
        sip_password: sipPassword,
      },
      asr: {
        url: asrUrl,
        token: asrToken,
      },
      tts: {
        url: ttsUrl,
        token: ttsToken,
      },
      id: clientId,
    }

    setIsLoading(true)
    const res = await startCampaign(payload)

    if (!res) {
      setCampaignStatus('failed')
      setIsLoading(false)
      return
    }

    setClientId(res)

    setCampaignStatus('running')
  }

  return (
    <Flex
      w="100%"
      h="calc(100vh - 98px)"
      px={'4.5rem'}
      justifyContent={'flex-start'}
      alignItems="center"
      as={'main'}
      gap="2rem"
      flexDirection={'column'}
    >
      <Flex w="100%" gap="3.75rem">
        <Flex w="50%" h="100%" flexDirection={'column'} gap={5}>
          <Settings {...assistantSettings} />
          <Settings {...asrSettings} />
          <Settings {...ttsSettings} />
          <Settings {...sipSettings} />
        </Flex>
        <Flex w="50%" h="100%" flexDirection={'column'}>
          <Contacts value={rawContacts} setter={setRawContacts} disabled={disabled} />
        </Flex>
      </Flex>
      <Flex
        w="100%"
        justifyContent={'flex-end'}
        alignItems={'center'}
        gap={'2.5rem'}
        color={colors[4]}
        fontWeight={400}
        fontSize={'18px'}
        lineHeight={'24px'}
      >
        {campaignStatus === 'failed' && <Text color={colors[1]}>возникла ошибка!</Text>}
        <Button
          onClick={changeStatus}
          disabled={isIncorrectSettings}
          fontSize={'18px'}
          variant={campaignStatus === 'running' ? 'primary' : 'secondary'}
          isLoading={isLoading}
          loadingText={buttonText}
          spinnerPlacement="end"
          _loading={{ opacity: 0.3, cursor: 'not-allowed' }}
        >
          {buttonText}
        </Button>
      </Flex>
    </Flex>
  )
})

Campaign.displayName = 'Campaign'

export { Campaign }
