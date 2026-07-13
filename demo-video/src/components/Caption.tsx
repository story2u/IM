import { interpolate, useCurrentFrame } from 'remotion'

export function Caption({ children }: { children: string }) {
  const frame = useCurrentFrame()
  const opacity = interpolate(frame, [0, 18], [0, 1], { extrapolateRight: 'clamp' })
  return <div style={{position:'absolute',left:180,right:180,bottom:58,display:'flex',justifyContent:'center',opacity}}><div style={{maxWidth:1250,padding:'17px 28px',borderRadius:8,background:'rgba(20,20,26,.9)',color:'white',fontSize:34,lineHeight:1.45,textAlign:'center'}}>{children}</div></div>
}
