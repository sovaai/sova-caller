import { Flex } from '@chakra-ui/react'
import React, { memo } from 'react'
import { Campaign } from './components/campaign/campaign'

import { Header } from './components/header/header'

const App = memo(() => (
  <Flex h="100vh" w="100vw" flexDirection={'column'} justifyContent="center" alignItems={'center'}>
    <Flex maxW={'1440px'} flexDirection={'column'} px={'3.75rem'} width={'100%'}>
      <Header />
      <Campaign />
    </Flex>
  </Flex>
))

App.displayName = 'App'

export { App }
