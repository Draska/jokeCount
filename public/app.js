new Vue({
    el: '#app',

    data: {
        jokers: [],
        jokersPlace: '',
        newJoker: '',
        joker: ''
    },

    created: function() {
        var self = this;
        this.score()
    },

    methods: {
        addJoke: function () {
            if(this.joker == '') {
                Materialize.toast("Who is the joker?", 2000)
                return
            }
            fetch(`/joke/${this.joker}`)
            this.joker = ''; // Reset field
            this.score()
        },

        addJoker: function () {
            if(this.newJoker == '') {
                Materialize.toast("Add a joker!", 2000)
                return
            }
            fetch(`/add/${this.newJoker}`)
            this.newJoker = ''; // Reset field
            this.score()
        },

        addThisJoke: function (joker) {
            if(joker == '') {
                Materialize.toast("Add a joker!", 2000)
                return
            }
            fetch(`/joke/${joker}`)
            //this.newJoker = ''; // Reset field
            this.score()
            location.reload()
        },

        score: function() {
            this.jokersPlace = ''
            fetch(`/score`)
            .then( response => {
                if(response.status !== 200) {
                    Materialize.toast("Failed retrieving jokes score!", 2000)
                    console.log('Whoops! Not the expected status! Status:' + response.status);
                    return
                }
                response.json()
                .then( data => {
                    console.log(data);
                    this.jokers = data;
                    this.jokers.forEach(joker => {
                        this.parseJoker(joker);
                    });
                })
            })
            .catch(err => {
                console.log('Error retrieving jokers: -S', err)
            });
        },

        parseJoker: function(joker) {
            this.jokersPlace += `<a href="#!" class="collection-item"><span class="badge">${joker.jokes}</span>${ emojione.toImage(joker.name)}</a>`
        }
    }
});
