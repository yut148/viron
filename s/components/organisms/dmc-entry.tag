dmc-entry.Entry
  .Entry__title 新しい管理画面を<br />作成する
  .Entry__form
    dmc-input(text="{ endpointURL }" placeholder="エンドポイントURL" pattern="https?://[\w/:%#\$&\?\(\)~\.=\+\-]+" onTextChange="{ handleEndpointURLChange }")
    dmc-textarea(text="{ memo }" placeholder="Writing..." maxlength="20" onTextChange="{ handleMemoChange }")
  .Entry__controls
    dmc-button(type="primary" onClick="{ handleRegisterButtonClick }" label="新規作成")
    dmc-button(type="secondary" onClick="{ handleCancelButtonClick }" label="キャンセル")
  script.
    import constants from '../../core/constants';
    import '../atoms/dmc-button.tag';
    import '../atoms/dmc-input.tag';
    import '../atoms/dmc-textarea.tag';

    const store = this.riotx.get();

    this.endpointURL = 'http://127.0.0.1:3000/swagger.json';
    this.memo = '';

    closeModal() {
      if (this.opts.isModal) {
        this.opts.modalCloser();
      }
    }

    handleEndpointURLChange(endpointURL) {
      this.endpointURL = endpointURL;
      this.update();
    }

    handleMemoChange(memo) {
      this.memo = memo;
      this.update();
    }

    // TODO: 上書きの場合は、そもそも登録ボタンを押せなくする
    handleRegisterButtonClick() {
      store.action(constants.ACTION_ENDPOINTS_ADD, this.endpointURL, this.memo)
      .then(() => {
        this.closeModal();
      }).catch(err => {
        store.action(constants.ACTION_TOAST_SHOW, {
          message: err.message
        });
      })
    }

    handleCancelButtonClick() {
      this.closeModal();
    }
