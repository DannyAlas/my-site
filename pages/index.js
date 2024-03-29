import NextLink from 'next/link'
import {
  Link,
  Container,
  Heading,
  Box,
  SimpleGrid,
  Button,
  List,
  ListItem,
  UnorderedList,
  useColorModeValue,
  chakra
} from '@chakra-ui/react'
import { ChevronRightIcon } from '@chakra-ui/icons'
import Layout from '../components/layouts/article'
import Section from '../components/section'
import { GridItem } from '../components/grid-item'
import {
  IoLogoGithub,
  IoLogoLinkedin,
  IoLogoYoutube,
  IoLogOut
} from 'react-icons/io5'
import Image from 'next/image'
import { compareDesc } from 'date-fns'
import { allPosts } from 'contentlayer/generated'
import { WorkGridItemPosts } from '../components/grid-item'

const ProfileImage = chakra(Image, {
  
  shouldForwardProp: prop => ['width', 'height', 'src', 'alt', 'loader' ].includes(prop)

})

export async function getServerSideProps({ req }) {
  const posts = allPosts.sort((a, b) => {
    return compareDesc(new Date(a.date), new Date(b.date))
  })
  return { props: { posts: posts, cookies: req.headers.cookie ?? '' } }
}

export default function Home({ posts }) {
  return (
    <Layout>
      <Container>
        <Box
          borderRadius="lg"
          mb={6}
          p={3}
          textAlign="center"
          bg={useColorModeValue('whiteAlpha.500', 'whiteAlpha.200')}
          css={{ backdropFilter: 'blur(10px)' }}
        >
          An inquisitive introspective introvert based in Boulder, Colorado.
        </Box>

        <Box display={{ md: 'flex' }}>
          <Box flexGrow={1}>
            <Heading as="h2" variant="page-title">
              Daniel Alas
            </Heading>
            <p>Curious ( Researcher / Developer / Designer )</p>
          </Box>
          <Box
            flexShrink={0}
            mt={{ base: 4, md: 0 }}
            ml={{ md: 6 }}
            textAlign="center"
          >
            <Box
              borderColor="whiteAlpha.800"
              borderWidth={2}
              borderStyle="solid"
              w="100px"
              h="100px"
              display="inline-block"
              borderRadius="full"
              overflow="hidden"
            >
              <ProfileImage
                loader={({ src }) => {
                  return src
                }}
                src="https://i.danielalas.com/01355563.png"
                alt="Profile image"
                // borderRadius="full"
                width="100%"
                height="100%"
              />
            </Box>
          </Box>
        </Box>

        <Section delay={0.1}>
          <Box align="center" my={4}>
            <NextLink href="/projects" passHref scroll={false}>
              <Button rightIcon={<ChevronRightIcon />} colorScheme="teal">
                Projects
              </Button>
            </NextLink>
          </Box>
        </Section>
        <Section delay={0.2}>
          <Heading as="h3" variant="section-title">
            Publications
          </Heading>
          <UnorderedList>
            <ListItem>
            <Link
                href="https://www.biorxiv.org/content/10.1101/2023.04.06.535959"
                target="_blank"
              >
              (2023) Monosynaptic inputs to ventral tegmental area glutamate and GABA co-transmitting neurons
              </Link>
            </ListItem>
          </UnorderedList>
        </Section>
        <Section delay={0.3}>
          <Heading as="h3" variant="section-title">
            Posts
          </Heading>
          <SimpleGrid columns={[1, 2, 2]} gap={6}>
            {posts.map(
              (post, idx) =>
                // if idx is less the 4 then show the post
                idx < 2 && (
                  <WorkGridItemPosts key={idx} post={post}></WorkGridItemPosts>
                )
            )}
          </SimpleGrid>
          <Box align="center" my={4}>
            <NextLink href="/posts" passHref scroll={false}>
              <Button rightIcon={<ChevronRightIcon />} colorScheme="teal">
                Popular posts
              </Button>
            </NextLink>
          </Box>
        </Section>
        <Section delay={0.3}>
          <Heading as="h3" variant="section-title">
            Notable Projects
          </Heading>
          <SimpleGrid columns={[1, 2, 2]} gap={6}>
            <GridItem
              href="/projects/Custom_NN"
              title="Custom Neurul Network"
              thumbnail="https://i.danielalas.com/40026d60"
            >
              Custom Neural Networks in Python
            </GridItem>
            <GridItem
              href="https://www.youtube.com/watch?v=I9IjFHPDNKw&t=143s"
              title="Free Agent Simulation"
              thumbnail="https://i.danielalas.com/810e452e"
            >
              Physarum Transport simulation built in Unity
            </GridItem>
          </SimpleGrid>
        </Section>
        <Section delay={0.3}>
          <Heading as="h3" variant="section-title">
            Find Me
          </Heading>
          <List>
            <SimpleGrid columns={4} gap={2}>
              <ListItem>
                <Link href="https://github.com/DannyAlas" target="_blank">
                  <Button
                    variant="ghost"
                    colorScheme="teal"
                    leftIcon={<IoLogoGithub />}
                  >
                    GitHub
                  </Button>
                </Link>
              </ListItem>
              <ListItem>
                <Link
                  href="https://www.linkedin.com/in/daniel-alas/"
                  target="_blank"
                >
                  <Button
                    variant="ghost"
                    colorScheme="teal"
                    leftIcon={<IoLogoLinkedin />}
                  >
                    Linkedin
                  </Button>
                </Link>
              </ListItem>
              <ListItem>
                <Link
                  href="https://www.youtube.com/channel/UCWdoRUs7EwO1vxkekbkCSxw"
                  target="_blank"
                >
                  <Button
                    variant="ghost"
                    colorScheme="teal"
                    leftIcon={<IoLogoYoutube />}
                  >
                    YouTube
                  </Button>
                </Link>
              </ListItem>
              <ListItem>
                <Link href="mailto:hi@danielalas.com" target="_blank">
                  <Button
                    variant="ghost"
                    colorScheme="teal"
                    leftIcon={<IoLogOut />}
                  >
                    Contact
                  </Button>
                </Link>
              </ListItem>
            </SimpleGrid>
          </List>
        </Section>
      </Container>
    </Layout>
  )
}
