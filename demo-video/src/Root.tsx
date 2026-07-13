import { Composition } from 'remotion'
import { ProductDemo } from './ProductDemo'

export function Root() {
  return <Composition id="ProductDemo" component={ProductDemo} durationInFrames={4500} fps={30} width={1920} height={1080} />
}
