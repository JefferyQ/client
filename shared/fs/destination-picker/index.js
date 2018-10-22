// @flow
import * as I from 'immutable'
import * as React from 'react'
import * as Types from '../../constants/types/fs'
import * as Constants from '../../constants/fs'
import * as Styles from '../../styles'
import * as Kb from '../../common-adapters'
import {isMobile} from '../../constants/platform'
import Rows from '../row/rows-container'
import * as F from '../common'
import Breadcrumb from '../header/breadcrumb-container'

type Props = {
  path: Types.Path,
  targetName: string,
  targetIconSpec: Types.PathItemIconSpec,
  onNewFolder: string,
}

const DestinationPicker = (props: Props) => (
  <Kb.Box style={styles.container}>
    <Kb.Box2 direction="horizontal" centerChildren={true} style={styles.header}>
      <Kb.Text type="Header">Move or Copy “</Kb.Text>
      <F.PathItemIcon small={true} style={styles.icon} spec={props.targetIconSpec} />
      <Kb.Text type="Header">{props.targetName}”</Kb.Text>
    </Kb.Box2>
    <Kb.Box style={styles.anotherHeader}>
      <Breadcrumb path={props.path} />
      <Kb.ClickableBox style={styles.newFolderBox} onClick={props.onNewFolder}>
        <Kb.Icon type="iconfont-folder-new" color={Styles.globalColors.blue} />
        <Kb.Text type="BodySemibold" style={styles.newFolderText}>
          Create new folder
        </Kb.Text>
      </Kb.ClickableBox>
    </Kb.Box>
    <Kb.Divider />
    <Kb.Box style={styles.rowsContainer}>
      <Rows path={props.path} sortSetting={Constants.sortByNameAsc} />
    </Kb.Box>
    <Kb.Box style={styles.footer}>
      <Kb.Button type="Secondary" label="Cancel" onClick={props.onCancel} />
    </Kb.Box>
  </Kb.Box>
)

const styles = Styles.styleSheetCreate({
  container: Styles.platformStyles({
    common: {
      ...Styles.globalStyles.flexBoxColumn,
    },
    isElectron: {
      width: 560,
      height: 480,
    },
  }),
  header: {
    marginTop: 32,
    marginBottom: 10,
  },
  anotherHeader: {
    ...Styles.globalStyles.flexBoxRow,
    height: 48,
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  newFolderBox: {
    ...Styles.globalStyles.flexBoxRow,
    height: 48,
    alignItems: 'center',
    padding: Styles.globalMargins.small,
  },
  newFolderText: {
    marginLeft: Styles.globalMargins.tiny,
    color: Styles.globalColors.blue,
  },
  icon: {
    marginRight: Styles.globalMargins.tiny,
    marginLeft: 2,
  },
  rowsContainer: {
    ...Styles.globalStyles.flexBoxColumn,
    ...Styles.globalStyles.fullHeight,
    flex: 1,
  },
  footer: {
    ...Styles.globalStyles.flexBoxRow,
    justifyContent: 'center',
    alignItems: 'center',
    padding: Styles.globalMargins.small,
  },
})

export default Kb.HeaderOrPopup(DestinationPicker)
