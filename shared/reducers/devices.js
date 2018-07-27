// @flow
import * as I from 'immutable'
import * as Constants from '../constants/devices'
import * as Types from '../constants/types/devices'
import * as DevicesGen from '../actions/devices-gen'

const initialState: Types.State = Constants.makeState()

export default function(state: Types.State = initialState, action: DevicesGen.Actions) {
  switch (action.type) {
    case DevicesGen.resetStore:
      return initialState
    case DevicesGen.loaded:
      console.error('aaa reducer loaded', action)
      return state.set('deviceMap', I.Map(action.payload.devices.map(d => [d.deviceID, d])))
    case DevicesGen.endangeredTLFsLoaded:
      return state.setIn(['endangeredTLFMap', action.payload.deviceID], I.Set(action.payload.tlfs))
    case DevicesGen.showRevokePage:
    case DevicesGen.showDevicePage: // fallthrough
      return state.set('selectedDeviceID', action.payload.deviceID)
    case DevicesGen.showPaperKeyPage:
      return state.set('newPaperkey', initialState.newPaperkey)
    case DevicesGen.paperKeyCreated:
      return state.set('newPaperkey', action.payload.paperKey)
    // Saga only actions
    case DevicesGen.revoke:
    case DevicesGen.revoked:
    case DevicesGen.load:
      return state
    default:
      /*::
      declare var ifFlowErrorsHereItsCauseYouDidntHandleAllActionTypesAbove: (action: empty) => any
      ifFlowErrorsHereItsCauseYouDidntHandleAllActionTypesAbove(action);
      */
      return state
  }
}
