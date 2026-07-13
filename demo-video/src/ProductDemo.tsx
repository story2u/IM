import { AbsoluteFill, Audio, Img, Loop, OffthreadVideo, Sequence, staticFile } from 'remotion'
import { BrandMark } from './components/BrandMark'
import { BrowserFrame } from './components/BrowserFrame'
import { Caption } from './components/Caption'
import { ChapterTitle } from './components/ChapterTitle'
import { zhCN } from './content/script.zh-CN'

const bg = '#f7f7fa'
const Screenshot = ({ src }: { src: string }) => <Img src={staticFile(src)} style={{width:'100%',height:'100%',objectFit:'cover',objectPosition:'top'}} />

export function ProductDemo() {
  return <AbsoluteFill style={{background:bg,fontFamily:'Inter, Arial, sans-serif',color:'#25242d'}}>
    <Audio src={staticFile('audio/background.wav')} volume={0.12} />
    <Sequence durationInFrames={180}><AbsoluteFill style={{display:'flex',alignItems:'center',justifyContent:'center',background:'#202027',color:'white'}}><div style={{textAlign:'center'}}><BrandMark /><h1 style={{fontSize:70,margin:'48px 0 18px'}}>不错过每一条藏在聊天里的商机</h1><p style={{fontSize:28,opacity:.68}}>{zhCN.intro}</p></div></AbsoluteFill></Sequence>
    <Sequence from={180} durationInFrames={900}><AbsoluteFill><ChapterTitle eyebrow="01 · 产品看板" title="从多 IM 消息中发现真实需求" /><BrowserFrame><Screenshot src="screenshots/web-home-desktop.png" /></BrowserFrame><Caption>{zhCN.dashboard}</Caption></AbsoluteFill></Sequence>
    <Sequence from={1080} durationInFrames={1800}><AbsoluteFill><ChapterTitle eyebrow="02 · 真实操作" title="识别、审核、生成草稿与筛选" /><BrowserFrame><Loop durationInFrames={900}><OffthreadVideo src={staticFile('raw/web-demo.webm')} muted style={{width:'100%',height:'100%',objectFit:'cover'}} /></Loop></BrowserFrame><Caption>{zhCN.workflow}</Caption></AbsoluteFill></Sequence>
    <Sequence from={2880} durationInFrames={720}><AbsoluteFill><ChapterTitle eyebrow="03 · Pi Agent" title="上下文、风险、联系人和行动建议" /><BrowserFrame><Screenshot src="screenshots/opportunity-detail.png" /></BrowserFrame><Caption>{zhCN.agent}</Caption></AbsoluteFill></Sequence>
    <Sequence from={3600} durationInFrames={600}><AbsoluteFill style={{padding:'150px 160px'}}><ChapterTitle eyebrow="04 · 多端协同" title="Web 完整运营，移动端及时处理" /><div style={{display:'grid',gridTemplateColumns:'1.2fr .8fr .8fr',gap:28,marginTop:110,height:590}}><div style={{border:'1px solid #ddd',borderRadius:10,overflow:'hidden'}}><Screenshot src="screenshots/dashboard-desktop.png" /></div>{['iOS App','Android App'].map((x)=><div key={x} style={{border:'1px solid #ddd',borderRadius:10,display:'flex',alignItems:'center',justifyContent:'center',background:'white',textAlign:'center'}}><div><div style={{fontSize:70}}>◎</div><h3 style={{fontSize:32}}>{x}</h3><p style={{fontSize:20,color:'#777'}}>Beta · 平台截图待生成</p></div></div>)}</div><Caption>{zhCN.mobile}</Caption></AbsoluteFill></Sequence>
    <Sequence from={4200} durationInFrames={300}><AbsoluteFill style={{display:'flex',alignItems:'center',justifyContent:'center',background:'#202027',color:'white'}}><div style={{textAlign:'center'}}><BrandMark /><h2 style={{fontSize:64,margin:'52px 0 22px'}}>{zhCN.outro}</h2><p style={{fontSize:28,opacity:.68}}>github.com/story2u/IM · im.story2u.xyz</p><p style={{fontSize:24,color:'#8d82ff',marginTop:28}}>AI 分析，人工决策</p></div></AbsoluteFill></Sequence>
  </AbsoluteFill>
}
